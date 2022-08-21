// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package lookup

import (
	"context"
	"strings"
	"sync/atomic"

	"github.com/apex/log"
	"github.com/lrstanley/geoip/internal/cache"
	"github.com/lrstanley/geoip/internal/dns"
	"github.com/lrstanley/geoip/internal/models"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

type Service struct {
	ctx    context.Context
	logger log.Interface
	config models.ConfigDB

	cache    *cache.Cache[string, *models.Response]
	metadata atomic.Pointer[maxminddb.Metadata]

	rslv *dns.Resolver
}

func NewService(ctx context.Context, logger log.Interface, config models.ConfigDB, rslv *dns.Resolver) *Service {
	return &Service{
		ctx:    ctx,
		logger: logger.WithFields(log.Fields{"src": "lookup"}),
		config: config,
		cache:  cache.New[string, *models.Response](config.CacheSize, config.CacheExpire),
		rslv:   rslv,
	}
}

func (s *Service) Metadata() *maxminddb.Metadata {
	return s.metadata.Load()
}

func (s *Service) MatchLanguage(lang string) (match string) {
	if lang == "" {
		return ""
	}

	supported := s.Metadata().Languages

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

	return ""
}
