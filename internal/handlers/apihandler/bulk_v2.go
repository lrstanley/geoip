// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package apihandler

import (
	"net/http"
	"sync"

	"github.com/lrstanley/chix"
	"github.com/lrstanley/geoip/internal/httpware"
	"github.com/lrstanley/geoip/internal/models"
	"github.com/sourcegraph/conc/pool"
	"golang.org/x/exp/slices"
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

	// Max 5 concurrent lookups from this request.
	pool := pool.New().WithMaxGoroutines(5)

	resp := &BulkResponse{
		Results: make([]*models.Response, 0, len(opts.Addresses)),
		Errors:  []*BulkError{},
	}

	for _, addr := range opts.Addresses {
		_, _, _, ok, err := h.limiter.Store.Take(r.Context(), h.limiter.Key(r))
		if err != nil {
			resp.AddError(addr, err)
			continue
		}

		if !ok {
			resp.AddError(addr, &models.ErrRateLimitExceeded{Address: addr})
			continue
		}

		pool.Go(func() {
			result, err := h.lookupSvc.Lookup(r.Context(), addr, &opts.LookupOptions)
			if err != nil {
				resp.AddError(addr, err)
				return
			}

			resp.AddResult(result)
		})
	}

	// Don't need to check ctx here, because we pass through the ctx to all goroutines
	// and the ctx is cancelled when the request is cancelled.
	pool.Wait()
	chix.JSON(w, r, http.StatusOK, resp)
}
