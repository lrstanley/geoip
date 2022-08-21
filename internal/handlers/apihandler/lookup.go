// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package apihandler

import (
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/chix"
	"github.com/lrstanley/geoip/internal/httpware"
	"github.com/lrstanley/geoip/internal/models"
)

func (h *handler) getLookup(w http.ResponseWriter, r *http.Request) {
	addr := strings.TrimSpace(chi.URLParam(r, "addr"))
	logger := chix.Log(r).WithField("lookup_addr", addr)

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

	req := &models.LookupRequest{
		Address:  addr,
		Language: httpware.GetLanguage(r),
	}

	result, err := h.lookupSvc.Lookup(r.Context(), req)
	if err != nil {
		logger.WithError(err).Error("error looking up addr")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if result.Cached {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
	}

	chix.JSON(w, r, http.StatusOK, result)
}
