// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package apihandler

import (
	"net/http"

	"github.com/lrstanley/chix"
	"github.com/lrstanley/geoip/internal/models"
)

func (h *handler) getMetadataV2(w http.ResponseWriter, r *http.Request) {
	metadata := chix.M{}

	for _, db := range models.Databases {
		m, ok := h.lookupSvc.Metadata.Load(db)
		if ok {
			metadata[db] = m
		} else {
			metadata[db] = chix.M{
				"DatabaseType": "unavailable",
			}
		}
	}

	chix.JSON(w, r, http.StatusOK, metadata)
}
