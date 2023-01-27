// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package lookup

import (
	"context"
	"strings"

	"github.com/apex/log"
	"github.com/lrstanley/geoip/internal/cache"
	"github.com/lrstanley/geoip/internal/dns"
	"github.com/lrstanley/geoip/internal/models"
	maxminddb "github.com/oschwald/maxminddb-golang"
	"github.com/puzpuzpuz/xsync"
)

type Service struct {
	ctx    context.Context
	logger log.Interface
	config models.ConfigDB

	geoCache *cache.Cache[string, *models.GeoQuery]
	asnCache *cache.Cache[string, *models.ASNQuery]
	Metadata *xsync.MapOf[string, *maxminddb.Metadata]

	rslv *dns.Resolver
}

func NewService(ctx context.Context, logger log.Interface, config models.ConfigDB, rslv *dns.Resolver) *Service {
	return &Service{
		ctx:      ctx,
		logger:   logger.WithFields(log.Fields{"src": "lookup"}),
		config:   config,
		geoCache: cache.New[string, *models.GeoQuery]("service_geo", config.CacheSize, config.CacheExpire),
		asnCache: cache.New[string, *models.ASNQuery]("service_asn", config.CacheSize, config.CacheExpire),
		Metadata: xsync.NewMapOf[*maxminddb.Metadata](),
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
		for i := 0; i < len(supported); i++ {
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
