// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package lookup

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/lrstanley/geoip/internal/dns"
	"github.com/lrstanley/geoip/internal/models"
	"github.com/lrstanley/x/sync/cache"
	"github.com/lrstanley/x/sync/cache/policy/lfu"
	"github.com/lrstanley/x/sync/conc"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

type Service struct {
	ctx    context.Context
	logger *slog.Logger
	config models.ConfigDB

	geoCache *cache.Cache[string, *models.GeoQuery]
	asnCache *cache.Cache[string, *models.ASNQuery]
	Metadata *conc.Map[string, *maxminddb.Metadata]

	rslv *dns.Resolver
}

func NewService(ctx context.Context, logger *slog.Logger, config models.ConfigDB, rslv *dns.Resolver) *Service {
	return &Service{
		ctx:    ctx,
		logger: logger.With("src", "lookup"),
		config: config,
		geoCache: cache.New(
			ctx,
			cache.WithJanitorInterval[string, *models.GeoQuery](1*time.Minute),
			cache.WithLFU[string, *models.GeoQuery](lfu.WithCapacity(config.CacheSize)),
			cache.WithDefaultEntryOptions[string, *models.GeoQuery](cache.WithExpiration(config.CacheExpire)),
		),
		asnCache: cache.New(
			ctx,
			cache.WithJanitorInterval[string, *models.ASNQuery](1*time.Minute),
			cache.WithLFU[string, *models.ASNQuery](lfu.WithCapacity(config.CacheSize)),
			cache.WithDefaultEntryOptions[string, *models.ASNQuery](cache.WithExpiration(config.CacheExpire)),
		),
		Metadata: &conc.Map[string, *maxminddb.Metadata]{},
		rslv:     rslv,
	}
}

func (s *Service) MatchLanguage(dbType string, languages []string) (match string) {
	languages = append(languages, s.config.DefaultLanguage, "en")

	var supported []string

	m, ok := s.Metadata.Load(dbType)
	if ok {
		supported = m.Languages
	} else {
		supported = append(supported, "en")
	}

	for _, lang := range languages {
		for i := range supported {
			if strings.EqualFold(lang, supported[i]) {
				return supported[i]
			}

			if j := strings.Index(supported[i], "-"); j > 0 {
				if strings.EqualFold(lang, supported[i][:j]) {
					return supported[i]
				}

				if k := strings.Index(lang, "-"); k > 0 {
					if strings.EqualFold(lang[:k], supported[i][:j]) {
						return supported[i]
					}
				}
			}

			if j := strings.Index(lang, "-"); j > 0 {
				if strings.EqualFold(lang[:j], supported[i]) {
					return supported[i]
				}
			}
		}
	}

	return ""
}
