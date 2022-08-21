// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package lookup

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/lrstanley/geoip/internal/models"
	"github.com/lrstanley/go-bogon"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

// Lookup does a geoip lookup of an address.
func (s *Service) Lookup(ctx context.Context, r *models.LookupRequest) (result *models.Response, err error) {
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

			return &models.Response{Error: fmt.Sprintf("invalid ip/host specified (or timed out): %s", r.Address)}, nil
		}
	}

	if is, _ := bogon.Is(ip.String()); is {
		return &models.Response{Error: "internal address"}, nil
	}

	geo, err := s.lookupGeo(ctx, ip)
	if err != nil {
		return nil, err
	}

	network, asn, err := s.lookupASN(ctx, ip)
	if err != nil {
		return nil, err
	}

	result = &models.Response{
		IP:               ip.String(),
		City:             geo.City.Names[r.Language],
		Country:          geo.Country.Names[r.Language],
		CountryCode:      geo.Country.Code,
		Continent:        geo.Continent.Names[r.Language],
		ContinentCode:    geo.Continent.Code,
		Lat:              geo.Location.Lat,
		Long:             geo.Location.Long,
		Timezone:         geo.Location.TimeZone,
		AccuracyRadiusKM: geo.Location.AccuracyRadiusKM,
		PostalCode:       geo.Postal.Code,
		ASN:              fmt.Sprintf("ASN%d", asn.AutonomousSystemNumber),
		ASNOrg:           asn.AutonomousSystemOrg,
		Network:          network.String(),
	}

	if v4 := ip.To4(); v4 != nil {
		result.IPType = 4
	} else if v6 := ip.To16(); v6 != nil {
		result.IPType = 6
	}

	var subdiv []string
	for i := 0; i < len(geo.Subdivisions); i++ {
		subdiv = append(subdiv, geo.Subdivisions[i].Names[r.Language])
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

func (s *Service) lookupASN(ctx context.Context, ip net.IP) (network *net.IPNet, query *models.ASNQuery, err error) {
	var db *maxminddb.Reader

	db, err = maxminddb.Open(s.config.ASNPath)
	if err != nil {
		s.logger.WithError(err).Error("error opening asn db")
		return nil, nil, err
	}
	defer db.Close()

	query = &models.ASNQuery{}
	if network, _, err = db.LookupNetwork(ip, &query); err != nil {
		s.logger.WithError(err).WithField("ip", ip.String()).Error("error looking up ip asn info")
		return nil, nil, err
	}

	return network, query, nil
}

func (s *Service) lookupGeo(ctx context.Context, ip net.IP) (query *models.GeoQuery, err error) {
	var db *maxminddb.Reader

	db, err = maxminddb.Open(s.config.GeoIPPath)
	if err != nil {
		s.logger.WithError(err).Error("error opening geoip db")
		return nil, err
	}
	defer db.Close()

	query = &models.GeoQuery{}
	if err = db.Lookup(ip, query); err != nil {
		s.logger.WithError(err).WithField("ip", ip.String()).Error("error looking up ip geoip info")
		return nil, err
	}

	return query, nil
}
