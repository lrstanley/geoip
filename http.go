// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"crypto/tls"
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
)

func initHTTP(closer chan struct{}) {
	r := chi.NewRouter()
	if flags.HTTP.Proxy {
		r.Use(middleware.RealIP)
	}

	r.Use(middleware.CloseNotify)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.GetHead)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.DefaultCompress)
	r.Use(dbDetailsMiddleware)

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

	limiterBackend := httprl.NewMap(10)
	limiter := &httprl.RateLimiter{
		Backend:  limiterBackend,
		Limit:    uint64(flags.HTTP.Limit),
		Interval: 60 * 60, // 1h.
		LimitExceededFunc: func(w http.ResponseWriter, r *http.Request) {
			debug.Printf(
				"connection %s has hit rate limit (limit: %s, reset: %s)",
				r.RemoteAddr,
				w.Header().Get("X-Ratelimit-Limit"),
				w.Header().Get("X-Ratelimit-Reset"),
			)
		},
		KeyMaker: httprl.DefaultKeyMaker, // This uses IP address by default.
	}
	limiterBackend.Start()
	defer limiterBackend.Stop()

	if flags.HTTP.Limit > 0 {
		r.With(cors.Handler, middleware.NoCache, limiter.Handle).Group(registerAPI)
	} else {
		r.With(cors.Handler, middleware.NoCache).Group(registerAPI)
	}

	srv := http.Server{
		Addr:         flags.HTTP.Bind,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if flags.HTTP.TLS.Use {
		srv.TLSConfig = &tls.Config{PreferServerCipherSuites: true}

		go func() {
			debug.Println("starting https server")
			err := srv.ListenAndServeTLS(flags.HTTP.TLS.Cert, flags.HTTP.TLS.Key)
			if err != nil {
				fmt.Printf("error in https server: %s\n", err)
				os.Exit(1)
			}
		}()
	} else {
		go func() {
			debug.Println("starting http server")
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
		debug.Printf("error while stopping http server: %s", err)
	}
}
