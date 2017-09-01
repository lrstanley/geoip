package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-web/httprl"
)

func initHTTP(closer chan struct{}) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.GetHead)

	if flags.HTTP.Proxy {
		r.Use(middleware.RealIP)
	}

	if flags.HTTP.CORS == nil || len(flags.HTTP.CORS) == 0 {
		flags.HTTP.CORS = []string{"*"}
	}
	cors := cors.New(cors.Options{
		AllowedOrigins: flags.HTTP.CORS,
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		ExposedHeaders: []string{"X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
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
				w.Header().Get("X-RateLimit-Limit"),
				w.Header().Get("X-RateLimit-Reset"),
			)
		},
		KeyMaker: httprl.DefaultKeyMaker, // This uses IP address by default.
	}
	limiterBackend.Start()
	defer limiterBackend.Stop()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK\n"))
	})

	if flags.HTTP.Limit > 0 {
		r.With(cors.Handler, limiter.Handle).Group(registerAPI)
	} else {
		r.With(cors.Handler).Group(registerAPI)
	}

	srv := http.Server{
		Addr:         flags.HTTP.Bind,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if flags.HTTP.TLS.Use {
		srv.TLSConfig = &tls.Config{PreferServerCipherSuites: true}

		go srv.ListenAndServeTLS(flags.HTTP.TLS.Cert, flags.HTTP.TLS.Key)
	} else {
		go srv.ListenAndServe()
	}

	<-closer
	fmt.Println("gracefully closing http connections")

	if err := srv.Close(); err != nil {
		debug.Printf("error while stopping http server: %s", err)
	}
}
