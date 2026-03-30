// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package apihandler

import (
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/geoip/internal/httpware"
	"github.com/lrstanley/geoip/internal/models"
)

func (h *handler) getLookupV1(w http.ResponseWriter, r *http.Request) {
	addr := strings.TrimSpace(chi.URLParam(r, "addr"))

	opts := &models.LookupOptions{}
	if err := chix.Bind(r, opts); err != nil {
		result := &models.Response{
			Query: addr,
			Error: err.Error(),
		}

		chix.JSON(w, r, http.StatusOK, result)
		return
	}

	if len(opts.Languages) == 0 {
		opts.Languages = httpware.GetLanguage(r)
	}

	// Allow users to query themselves without having to have them specify
	// their own IP address. Note that this will not work if you are querying
	// the IP address locally.
	if self := strings.ToLower(addr); self == "self" || self == "me" {
		if strings.Contains(r.RemoteAddr, ":") {
			addr, _, _ = net.SplitHostPort(r.RemoteAddr)
		} else {
			addr = r.RemoteAddr
		}
	}

	result, err := h.lookupSvc.Lookup(r.Context(), addr, opts)
	if err != nil {
		if models.IsClientError(err) {
			result = &models.Response{
				Query: addr,
				Error: err.Error(),
			}

			chix.JSON(w, r, http.StatusOK, result)
			return
		}

		chix.LogError(
			r.Context(), "error looking up addr",
			slog.String("lookup_addr", addr),
			slog.Any("error", err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	chix.JSON(w, r, http.StatusOK, result)
}
