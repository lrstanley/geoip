// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"embed"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/lrstanley/chix"
	"github.com/lrstanley/geoip/internal/handlers/apihandler"
	"github.com/lrstanley/geoip/internal/httpware"
	"github.com/lrstanley/geoip/internal/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:generate sh -c "mkdir -vp public/dist;touch public/dist/index.html"
//go:embed all:public/dist
var staticFS embed.FS

func httpServer(ctx context.Context) *http.Server {
	chix.AddErrorResolver(models.ErrorResolver)

	r := chi.NewRouter()

	limiter := httpware.NewLimiter(cli.Flags.HTTP, 1*time.Hour)

	if len(cli.Flags.HTTP.TrustedProxies) > 0 {
		r.Use(chix.UseRealIPCLIOpts(cli.Flags.HTTP.TrustedProxies))
	}

	// Core middeware.
	r.Use(
		chix.UseContextIP,
		middleware.RequestID,
		chix.UseStructuredLogger(logger),
		chix.UseIf(cli.Flags.HTTP.Metrics, chix.UsePrometheus),
		chix.UseDebug(cli.Debug),
		chix.Recoverer,
		middleware.Maybe(middleware.StripSlashes, func(r *http.Request) bool {
			return !strings.HasPrefix(r.URL.Path, "/debug/")
		}),
		middleware.Compress(5),
	)

	if cli.Flags.HTTP.MaxConcurrent > 0 {
		r.Use(middleware.Throttle(cli.Flags.HTTP.MaxConcurrent))
	}

	if cli.Flags.HTTP.CORS == nil || len(cli.Flags.HTTP.CORS) == 0 {
		cli.Flags.HTTP.CORS = []string{"*"}
	}

	// Security related.
	if !cli.Debug && cli.Flags.HTTP.HSTS {
		r.Use(middleware.SetHeader("Strict-Transport-Security", "max-age=31536000"))
	}
	r.Use(
		cors.New(cors.Options{
			AllowedOrigins: cli.Flags.HTTP.CORS,
			AllowedMethods: []string{"GET", "HEAD", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Content-Type"},
			ExposedHeaders: []string{
				"X-Maxmind-Type", "X-Maxmind-Version",
				"X-Ratelimit-Limit", "X-Ratelimit-Remaining", "X-Ratelimit-Reset",
				"X-Cache",
			},
			MaxAge: 3600,
		}).Handler,
		chix.UseHeaders(map[string]string{
			"Content-Security-Policy": "default-src 'self'; style-src 'self' 'unsafe-inline'; object-src 'none'; child-src 'none'; frame-src 'none'; worker-src 'none'; img-src 'self' data: https://*.openstreetmap.org https://hatscripts.github.io;",
			"X-Frame-Options":         "DENY",
			"X-Content-Type-Options":  "nosniff",
			"Referrer-Policy":         "no-referrer-when-downgrade",
			"Permissions-Policy":      "clipboard-write=(self)",
		}),
	)

	if cli.Debug {
		r.With(chix.UsePrivateIP).Mount("/debug", middleware.Profiler())
	}

	if cli.Flags.HTTP.Metrics {
		r.With(chix.UsePrivateIP).Mount("/metrics", promhttp.Handler())
	}

	r.Route("/api", apihandler.New(lookupSvc, limiter).Route)

	r.NotFound(chix.UseStatic(ctx, &chix.Static{
		FS:         staticFS,
		CatchAll:   true,
		AllowLocal: cli.Debug,
		Path:       "public/dist",
		SPA:        true,
		Headers: map[string]string{
			"Vary":          "Accept-Encoding",
			"Cache-Control": "public, max-age=7776000",
		},
	}).ServeHTTP)

	return &http.Server{
		Addr:         cli.Flags.HTTP.BindAddr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
