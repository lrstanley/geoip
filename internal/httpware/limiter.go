// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package httpware

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/lrstanley/chix"
	"github.com/lrstanley/geoip/internal/models"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
)

type contextKey string

const (
	// HeaderRateLimitLimit, HeaderRateLimitRemaining, and HeaderRateLimitReset
	// are the recommended return header values from IETF on rate limiting. Reset
	// is in UTC time.
	HeaderRateLimitLimit     = "X-RateLimit-Limit"
	HeaderRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderRateLimitReset     = "X-RateLimit-Reset"

	// HeaderRetryAfter is the header used to indicate when a client should retry
	// requests (when the rate limit expires), in UTC time.
	HeaderRetryAfter = "Retry-After"

	ctxLimiterKey contextKey = "limiter"
)

type Limiter struct {
	config models.ConfigHTTP
	Store  limiter.Store
}

func NewLimiter(config models.ConfigHTTP, window time.Duration) *Limiter {
	store, err := memorystore.New(&memorystore.Config{
		Tokens:        uint64(config.Limit),
		Interval:      window,
		SweepInterval: 10 * time.Minute,
		SweepMinTTL:   2 * time.Hour,
	})
	if err != nil {
		panic(err)
	}

	return &Limiter{
		config: config,
		Store:  store,
	}
}

func (l *Limiter) Key(r *http.Request) string {
	key, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil || key == "" {
		return r.RemoteAddr
	}
	return key
}

func (l *Limiter) Skip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxLimiterKey, true)))
	})
}

func (l *Limiter) IsSkipped(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	v := ctx.Value(ctxLimiterKey)
	if v == nil {
		return false
	}
	return v.(bool)
}

// Limit returns the HTTP handler as a middleware. Use Limiter.Skip() to skip
// the rate limiting logic, and only return the headers.
func (l *Limiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		skip := l.IsSkipped(ctx)

		key := l.Key(r)

		var limit, remaining, reset uint64
		var ok bool
		var err error

		if skip {
			// Just query the store.
			ok = true
			limit, remaining, err = l.Store.Get(ctx, key)
		} else {
			// Take from the store.
			limit, remaining, reset, ok, err = l.Store.Take(ctx, key)
		}

		if err != nil {
			chix.Error(w, r, chix.WrapCode(http.StatusInternalServerError))
			return
		}

		// Set headers (we do this regardless of whether the request is permitted).
		w.Header().Set(HeaderRateLimitLimit, strconv.FormatUint(limit, 10))
		w.Header().Set(HeaderRateLimitRemaining, strconv.FormatUint(remaining, 10))

		resetTime := time.Unix(0, int64(reset)).UTC().Format(time.RFC1123)
		if !skip {
			w.Header().Set(HeaderRateLimitReset, resetTime)
		}

		// Fail if there were no tokens remaining.
		if !ok {
			w.Header().Set(HeaderRetryAfter, resetTime)
			chix.Error(w, r, chix.WrapCode(http.StatusTooManyRequests))
			return
		}

		next.ServeHTTP(w, r)
	})
}
