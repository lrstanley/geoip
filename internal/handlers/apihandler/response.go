// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package apihandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lrstanley/chix"
	"github.com/lrstanley/geoip/internal/models"
)

func apiResponse(w http.ResponseWriter, r *http.Request, result *models.GeoResult, filters []string) {
	if result.Cached {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
	}

	var err error

	if len(filters) > 0 {
		if result.Error != "" {
			fmt.Fprintf(w, "err: %s", result.Error)
			return
		}

		base := make(map[string]*json.RawMessage)
		var tmp []byte

		tmp, err = json.Marshal(result)
		if err != nil {
			panic(err)
		}

		if err = json.Unmarshal(tmp, &base); err != nil {
			panic(err)
		}

		// TODO: wtf is this?
		out := make([]string, len(filters))
		for i := 0; i < len(filters); i++ {
			out[i] = strings.ReplaceAll(fmt.Sprintf("%s", *base[filters[i]]), "\"", "")
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strings.Join(out, "|")))
		return
	}

	chix.JSON(w, r, http.StatusOK, result)
}
