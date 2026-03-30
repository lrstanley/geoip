// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"embed"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/geoip/internal/handlers/apihandler"
	"github.com/lrstanley/geoip/internal/models"
)

//go:generate sh -c "mkdir -vp public/dist;touch public/dist/index.html"
//go:embed all:public/dist
var staticFS embed.FS

func httpServer(_ context.Context, logger *slog.Logger) *http.Server {
	cfg := chix.NewConfig().
		SetLogger(logger).
		SetAPIBasePath("/api").
		AddErrorResolvers(models.ErrorResolver)

	r := chi.NewRouter()

	rl := httprate.NewRateLimiter(
		cli.Flags.HTTP.Limit,
		time.Hour,
		httprate.WithKeyByIP(),
	)

	if len(cli.Flags.HTTP.TrustedProxies) > 0 {
		r.Use(chix.UseRealIPStringOpts(cli.Flags.HTTP.TrustedProxies))
	}

	// Core middleware.
	r.Use(
		cfg.Use(),
		chix.UseContextIP(),
		chix.UseRequestID(),
		chix.UseStructuredLogger(nil),
		chix.UseDebug(cli.Debug),
		chix.UseStripSlashes(),
		middleware.Compress(5),
	)

	if cli.Flags.HTTP.MaxConcurrent > 0 {
		r.Use(middleware.Throttle(cli.Flags.HTTP.MaxConcurrent))
	}

	if len(cli.Flags.HTTP.CORS) == 0 {
		cli.Flags.HTTP.CORS = []string{"*"}
	}

	// Security related.
	if !cli.Debug && cli.Flags.HTTP.HSTS {
		r.Use(middleware.SetHeader("Strict-Transport-Security", "max-age=31536000"))
	}
	r.Use(
		chix.UseCrossOriginResourceSharing(&chix.CORSConfig{
			AllowedOrigins: cli.Flags.HTTP.CORS,
			AllowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
			AllowedHeaders: []string{"Accept", "Content-Type"},
			ExposedHeaders: []string{
				"X-Maxmind-Type", "X-Maxmind-Version",
				"X-Ratelimit-Limit", "X-Ratelimit-Remaining", "X-Ratelimit-Reset",
				"X-Cache",
			},
			MaxAge: time.Hour,
		}),
		chix.UseHeaders(map[string]string{
			"Content-Security-Policy": "default-src 'self'; style-src 'self' 'unsafe-inline'; object-src 'none'; child-src 'none'; frame-src 'none'; worker-src 'none'; img-src 'self' data: https://*.openstreetmap.org https://hatscripts.github.io;",
			"X-Frame-Options":         "DENY",
			"X-Content-Type-Options":  "nosniff",
			"Referrer-Policy":         "no-referrer-when-downgrade",
			"Permissions-Policy":      "clipboard-write=(self)",
		}),
	)

	if cli.Debug {
		r.With(chix.UsePrivateIP()).Mount("/debug", middleware.Profiler())
	}

	r.Route("/api", apihandler.New(lookupSvc, rl).Route)

	staticHandler := chix.UseStatic(&chix.StaticConfig{
		FS:         staticFS,
		CatchAll:   true,
		AllowLocal: cli.Debug,
		Path:       "public/dist",
		SPA:        true,
	})
	r.NotFound(chix.UseHeaders(map[string]string{
		"Vary":          "Accept-Encoding",
		"Cache-Control": "public, max-age=7776000",
	})(staticHandler).ServeHTTP)

	return &http.Server{
		Addr:         cli.Flags.HTTP.BindAddr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
