// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package lookup

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync/atomic"

	"github.com/apex/log"
	"github.com/lrstanley/geoip/internal/cache"
	"github.com/lrstanley/geoip/internal/dns"
	"github.com/lrstanley/geoip/internal/models"
	"github.com/lrstanley/go-bogon"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

type Service struct {
	ctx    context.Context
	logger log.Interface
	config models.ConfigDB

	cache    *cache.Cache[string, *models.GeoResult]
	metadata atomic.Pointer[maxminddb.Metadata]

	rslv *dns.Resolver
}

func NewService(ctx context.Context, logger log.Interface, config models.ConfigDB, rslv *dns.Resolver) *Service {
	return &Service{
		ctx: ctx,
		logger: logger.WithFields(log.Fields{
			"src":  "lookup",
			"path": config.Path,
		}),
		config: config,
		cache:  cache.New[string, *models.GeoResult](config.CacheSize, config.CacheExpire),
		rslv:   rslv,
	}
}

// Lookup does a geoip lookup of an address.
func (s *Service) Lookup(ctx context.Context, r *models.LookupRequest) (*models.GeoResult, error) {
	var result *models.GeoResult
	var err error

	if r.Language == "" {
		r.Language = s.config.DefaultLanguage
	}

	if val := s.cache.Get(r.CacheID()); val != nil {
		val.Cached = true
		return val, nil
	}

	ip := net.ParseIP(r.Address)
	if ip == nil {
		ip, err = s.rslv.GetIP(ctx, r.Address)
		if err != nil || ip == nil {
			s.logger.WithError(err).WithField("addr", r.Address).Debug("error looking up addr as hostname")

			return &models.GeoResult{Error: fmt.Sprintf("invalid ip/host specified: %s", r.Address)}, nil
		}
	}

	if is, _ := bogon.Is(ip.String()); is {
		return &models.GeoResult{Error: "internal address"}, nil
	}

	db, err := maxminddb.Open(s.config.Path)
	if err != nil {
		s.logger.WithError(err).Error("error opening db")
		return nil, err
	}

	var query models.GeoQuery

	err = db.Lookup(ip, &query)
	db.Close()

	if err != nil {
		return nil, err
	}

	result = &models.GeoResult{
		IP:            ip,
		City:          query.City.Names[r.Language],
		Country:       query.Country.Names[r.Language],
		CountryCode:   query.Country.Code,
		Continent:     query.Continent.Names[r.Language],
		ContinentCode: query.Continent.Code,
		Lat:           query.Location.Lat,
		Long:          query.Location.Long,
		Timezone:      query.Location.TimeZone,
		PostalCode:    query.Postal.Code,
		Proxy:         query.Traits.Proxy,
	}

	var subdiv []string
	for i := 0; i < len(query.Subdivisions); i++ {
		subdiv = append(subdiv, query.Subdivisions[i].Names[r.Language])
	}
	result.Subdivision = strings.Join(subdiv, ", ")

	var summary []string
	if result.City != "" {
		summary = append(summary, result.City)
	}

	if result.Subdivision != "" && result.City != result.Subdivision {
		summary = append(summary, result.Subdivision)
	}

	if result.Country != "" && len(summary) == 0 {
		summary = append(summary, result.Country)
	} else if result.CountryCode != "" {
		summary = append(summary, result.CountryCode)
	}

	if result.Continent != "" && len(summary) == 0 {
		summary = append(summary, result.Continent)
	} else if result.ContinentCode != "" && result.Subdivision == "" && result.City == "" {
		summary = append(summary, result.ContinentCode)
	}

	result.Summary = strings.Join(summary, ", ")

	if result.Summary == "" {
		result.Error = "no results found"
	}

	result.Host, err = s.rslv.GetReverse(ctx, ip)
	if err != nil {
		s.logger.WithError(err).WithField("ip", ip.String()).Debug("error looking up reverse dns for ip")
	}

	s.cache.Set(r.CacheID(), result)
	return result, nil
}

func (s *Service) Metadata() *maxminddb.Metadata {
	return s.metadata.Load()
}
