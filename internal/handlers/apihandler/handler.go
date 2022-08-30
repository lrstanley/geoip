// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package apihandler

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lrstanley/chix"
	"github.com/lrstanley/geoip/internal/httpware"
	"github.com/lrstanley/geoip/internal/lookup"
)

type handler struct {
	lookupSvc *lookup.Service
	limiter   *httpware.Limiter
}

func New(lookupSvc *lookup.Service, limiter *httpware.Limiter) *handler {
	return &handler{
		lookupSvc: lookupSvc,
		limiter:   limiter,
	}
}

func (h *handler) Route(r chi.Router) {
	r.Use(
		middleware.NoCache,
		middleware.GetHead,
	)

	// v1 API.
	r.With(h.limiter.Limit).Get("/{addr}", h.getLookupV1)
	r.With(h.limiter.Skip, h.limiter.Limit).Get("/ping", h.ping)

	// v2 API.
	r.Get("/v2/openapi.yaml", h.getOpenAPI)
	r.With(h.limiter.Limit).Get("/v2/lookup/{addr}", h.getLookupV2)
	r.With(h.limiter.Limit).Post("/v2/bulk", h.postBulkV2)
	r.With(h.limiter.Limit).Get("/v2/metadata", h.getMetadataV2)
	r.With(h.limiter.Skip, h.limiter.Limit).Get("/v2/ping", h.ping)
	r.NotFound(h.notFound)
}

func (h *handler) notFound(w http.ResponseWriter, r *http.Request) {
	chix.Error(w, r, chix.WrapCode(http.StatusNotFound))
}

func (h *handler) ping(w http.ResponseWriter, r *http.Request) {
	chix.JSON(w, r, http.StatusOK, chix.M{"pong": true})
}

//go:embed openapi_v2.yaml
var openapiv2 string

func (h *handler) getOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.Write([]byte(openapiv2))
}
