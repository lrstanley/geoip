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
	"github.com/go-chi/httprate"
	"github.com/lrstanley/chix"
	"github.com/lrstanley/geoip/internal/handlers/apihandler"
	"github.com/lrstanley/geoip/internal/httpware"
)

//go:generate touch public/dist/index.html
//go:embed all:public/dist
var staticFS embed.FS

func httpServer(ctx context.Context) *http.Server {
	r := chi.NewRouter()

	if len(cli.Flags.HTTP.TrustedProxies) > 0 {
		r.Use(chix.UseRealIP(cli.Flags.HTTP.TrustedProxies, chix.OptUseXForwardedFor))
	}

	// Core middeware.
	r.Use(
		chix.UseContextIP,
		middleware.RequestID,
		chix.UseStructuredLogger(logger),
		chix.UseDebug(cli.Debug),
		middleware.Recoverer,
		middleware.Maybe(middleware.StripSlashes, func(r *http.Request) bool {
			return !strings.HasPrefix(r.URL.Path, "/debug/")
		}),
		middleware.Compress(5),
		httpware.UseMetadata(lookupSvc),
		httpware.UseLanguage(lookupSvc),
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
		// middleware.SetHeader(
		// 	"Content-Security-Policy",
		// 	"default-src 'self'; media-src * data:; style-src 'self' 'unsafe-inline'; object-src 'none'; child-src 'none'; frame-src 'none'; worker-src 'none'",
		// ),
		middleware.SetHeader("X-Frame-Options", "DENY"),
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("Referrer-Policy", "no-referrer-when-downgrade"),
		middleware.SetHeader("Permissions-Policy", "clipboard-write=(self)"),
		httprate.LimitByIP(cli.Flags.HTTP.Limit, 1*time.Hour),
	)

	if cli.Debug {
		r.With(chix.UsePrivateIP).Mount("/debug", middleware.Profiler())
	}

	r.Route("/api", apihandler.New(lookupSvc).Route)

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
