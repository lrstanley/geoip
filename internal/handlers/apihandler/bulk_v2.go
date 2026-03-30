// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package apihandler

import (
	"net/http"
	"slices"
	"sync"

	"github.com/go-chi/httprate"
	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/geoip/internal/httpware"
	"github.com/lrstanley/geoip/internal/models"
	"github.com/lrstanley/x/sync/conc"
)

type BulkRequest struct {
	models.LookupOptions

	Addresses []string `form:"addresses" json:"addresses" validate:"required,min=1,max=25,dive,hostname|ip|fqdn,required"`
}

type BulkError struct {
	Query string `json:"query"`
	Error string `json:"error"`
}

type BulkResponse struct {
	mu      sync.Mutex
	Results []*models.Response `json:"results"`
	Errors  []*BulkError       `json:"errors"`
}

func (b *BulkResponse) AddResult(result *models.Response) {
	b.mu.Lock()
	b.Results = append(b.Results, result)
	b.mu.Unlock()
}

func (b *BulkResponse) AddError(query string, err error) {
	b.mu.Lock()
	if models.IsClientError(err) {
		b.Errors = append(b.Errors, &BulkError{
			Query: query,
			Error: err.Error(),
		})
	} else {
		b.Errors = append(b.Errors, &BulkError{
			Query: query,
			Error: "Internal server error",
		})
	}
	b.mu.Unlock()
}

func (h *handler) postBulkV2(w http.ResponseWriter, r *http.Request) {
	opts := &BulkRequest{}

	// Disable host lookup by default for bulk requests, to reduce the amount
	// of dns requests needed, and improve response time.
	opts.LookupOptions.DisableHostLookup = true

	if err := chix.Bind(r, opts); err != nil {
		chix.Error(w, r, err)
		return
	}

	// Sort, and remove duplicates.
	slices.Sort(opts.Addresses)
	opts.Addresses = slices.Compact(opts.Addresses)

	if len(opts.Languages) == 0 {
		opts.Languages = httpware.GetLanguage(r)
	}

	key, err := httprate.KeyByIP(r)
	if err != nil {
		chix.Error(w, r, err)
		return
	}

	g := conc.NewGroup().WithMaxGoroutines(5)

	resp := &BulkResponse{
		Results: make([]*models.Response, 0, len(opts.Addresses)),
		Errors:  []*BulkError{},
	}

	for _, addr := range opts.Addresses {
		if h.limiter.OnLimit(w, r, key) {
			resp.AddError(addr, &models.ErrRateLimitExceeded{Address: addr})
			continue
		}

		g.Go(func() {
			result, lerr := h.lookupSvc.Lookup(r.Context(), addr, &opts.LookupOptions)
			if lerr != nil {
				resp.AddError(addr, lerr)
				return
			}

			resp.AddResult(result)
		})
	}

	g.Wait()
	chix.JSON(w, r, http.StatusOK, resp)
}
