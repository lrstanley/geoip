// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-web/httprl"
	"github.com/lrstanley/recoverer"
)

var apiPong = map[string]bool{
	"pong": true,
}

var mapLimiter = NewMapLimiter(10)

func initHTTP(closer chan struct{}) {
	r := chi.NewRouter()
	if flags.HTTP.Proxy {
		r.Use(middleware.RealIP)
	}

	r.Use(middleware.CloseNotify)
	r.Use(recoverer.New(recoverer.Options{Logger: os.Stderr, Show: flags.Debug, Simple: false}))
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.DefaultCompress)
	r.Use(dbDetailsMiddleware)

	if flags.Debug {
		r.Mount("/debug", middleware.Profiler())
	}

	r.Mount("/static/dist", http.StripPrefix("/static/dist", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		http.FileServer(rice.MustFindBox("public/dist").HTTPBox()).ServeHTTP(w, r)
	})))
	r.Mount("/static", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		http.FileServer(rice.MustFindBox("public/static").HTTPBox()).ServeHTTP(w, r)
	})))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}

		w.Write(rice.MustFindBox("public/static/html").MustBytes("index.html"))
	})

	if flags.HTTP.CORS == nil || len(flags.HTTP.CORS) == 0 {
		flags.HTTP.CORS = []string{"*"}
	}
	cors := cors.New(cors.Options{
		AllowedOrigins: flags.HTTP.CORS,
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		ExposedHeaders: []string{"X-Ratelimit-Limit", "X-Ratelimit-Remaining", "X-Ratelimit-Reset"},
		MaxAge:         3600,
	})

	limiter := &httprl.RateLimiter{
		Backend:  mapLimiter,
		Limit:    uint64(flags.HTTP.Limit),
		Interval: 60 * 60, // 1h.
		LimitExceededFunc: func(w http.ResponseWriter, r *http.Request) {
			logger.Printf(
				"connection %s has hit rate limit (limit: %s, reset: %s)",
				r.RemoteAddr,
				w.Header().Get("X-Ratelimit-Limit"),
				w.Header().Get("X-Ratelimit-Reset"),
			)
		},
		KeyMaker: httprl.DefaultKeyMaker, // This uses IP address by default.
	}

	mapLimiter.Start()
	defer mapLimiter.Stop()

	if flags.HTTP.Limit > 0 {
		r.With(cors.Handler, middleware.NoCache, limiter.Handle).Group(registerAPI)
	} else {
		r.With(cors.Handler, middleware.NoCache).Group(registerAPI)
	}

	// Register the /api/ping route separately, as it shouldn't be counted
	// towards API limits. This endpoint will both let users verify that the
	// service is functional, but also let them use headers to check API
	// limit information. This endpoint is the only one which has HTTP HEAD
	// support.
	r.With(middleware.NoCache, rateHeaderMiddleware).Get("/api/ping", pingHandler)
	r.With(middleware.NoCache, rateHeaderMiddleware).Head("/api/ping", pingHandler)

	srv := http.Server{
		Addr:         flags.HTTP.Bind,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if flags.HTTP.TLS.Use {
		srv.TLSConfig = &tls.Config{PreferServerCipherSuites: true}

		go func() {
			logger.Println("starting https server")
			err := srv.ListenAndServeTLS(flags.HTTP.TLS.Cert, flags.HTTP.TLS.Key)
			if err != nil {
				fmt.Printf("error in https server: %s\n", err)
				os.Exit(1)
			}
		}()
	} else {
		go func() {
			logger.Println("starting http server")
			err := srv.ListenAndServe()
			if err != nil {
				fmt.Printf("error in http server: %s\n", err)
				os.Exit(1)
			}
		}()
	}

	<-closer
	fmt.Println("gracefully closing http connections")

	if err := srv.Close(); err != nil {
		logger.Printf("error while stopping http server: %s", err)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = enc.Encode(apiPong)
}
