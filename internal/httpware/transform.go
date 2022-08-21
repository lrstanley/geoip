// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package httpware

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/lrstanley/geoip/internal/lookup"
)

type ctxKey string

const ctxLanguage ctxKey = "language"

var reLanguage = regexp.MustCompile(`^[^a-zA-Z, ]+.*?$`)

func UseLanguage(lookupSvc *lookup.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prioritize "lang" query param.
			lang := lookupSvc.MatchLanguage(strings.ReplaceAll(r.FormValue("lang"), " ", ""))

			if lang == "" {
				// Try to get the language from the Accept-Language header.
				for _, l := range strings.Split(reLanguage.ReplaceAllString(r.Header.Get("Accept-Language"), ""), ",") {
					lang = lookupSvc.MatchLanguage(l)
					if lang != "" {
						break
					}
				}
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxLanguage, lang)))
		})
	}
}

func GetLanguage(r *http.Request) string {
	return r.Context().Value(ctxLanguage).(string)
}
