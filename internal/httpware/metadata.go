// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package httpware

import (
	"fmt"
	"net/http"

	"github.com/lrstanley/geoip/internal/lookup"
)

func UseMetadata(lookupSvc *lookup.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m := lookupSvc.Metadata()
			if m == nil {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("X-Maxmind-Build", fmt.Sprintf("%d-%d", m.IPVersion, m.BuildEpoch))
			w.Header().Set("X-Maxmind-Type", m.DatabaseType)

			next.ServeHTTP(w, r)
		})
	}
}
