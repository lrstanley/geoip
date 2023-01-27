// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package lookup

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/lrstanley/geoip/internal/metrics"
	"github.com/lrstanley/geoip/internal/models"
	"github.com/lrstanley/go-bogon"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

// Lookup does a geoip lookup of an address.
func (s *Service) Lookup(ctx context.Context, addr string, r *models.LookupOptions) (result *models.Response, err error) {
	langGeo := s.MatchLanguage(models.DatabaseGeoIP, r.Languages)

	ip := net.ParseIP(addr)
	if ip == nil {
		ip, err = s.rslv.GetIP(ctx, addr)
		if err != nil || ip == nil {
			return nil, &models.ErrHostResolve{Address: addr, Err: err}
		}
	}

	if is, _ := bogon.Is(ip.String()); is {
		return nil, &models.ErrInternalAddress{Address: addr}
	}

	geo, err := s.lookupGeo(ctx, ip)
	if err != nil {
		return nil, err
	}

	asn, err := s.lookupASN(ctx, ip)
	if err != nil {
		return nil, err
	}

	result = &models.Response{
		Query:            addr,
		IP:               ip.String(),
		City:             geo.City.Names[langGeo],
		Country:          geo.Country.Names[langGeo],
		CountryCode:      geo.Country.Code,
		Continent:        geo.Continent.Names[langGeo],
		ContinentCode:    geo.Continent.Code,
		Lat:              geo.Location.Lat,
		Long:             geo.Location.Long,
		Timezone:         geo.Location.TimeZone,
		AccuracyRadiusKM: geo.Location.AccuracyRadiusKM,
		PostalCode:       geo.Postal.Code,
	}

	if asn.AutonomousSystemNumber > 0 {
		result.ASN = fmt.Sprintf("AS%d", asn.AutonomousSystemNumber)
		result.ASNOrg = asn.AutonomousSystemOrg
	}

	if asn.Network != nil {
		result.Network = asn.Network.String()
	}

	if v4 := ip.To4(); v4 != nil {
		result.IPType = 4
	} else if v6 := ip.To16(); v6 != nil {
		result.IPType = 6
	}

	var subdiv []string
	for i := 0; i < len(geo.Subdivisions); i++ {
		subdiv = append(subdiv, geo.Subdivisions[i].Names[langGeo])
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
		return nil, &models.ErrNotFound{Address: addr}
	}

	if !r.DisableHostLookup {
		result.Host, err = s.rslv.GetReverse(ctx, ip)
		if err != nil {
			s.logger.WithError(err).WithField("ip", ip.String()).Debug("error looking up reverse dns for ip")
		}
	}

	return result, nil
}

func (s *Service) lookupASN(ctx context.Context, ip net.IP) (query *models.ASNQuery, err error) {
	metrics.LookupCount.WithLabelValues("asn").Inc()

	if val := s.asnCache.Get(ip.String()); val != nil {
		return val, nil
	}

	var db *maxminddb.Reader

	db, err = maxminddb.Open(s.config.ASNPath)
	if err != nil {
		s.logger.WithError(err).Error("error opening asn db")
		return nil, err
	}
	defer db.Close()

	query = &models.ASNQuery{}
	if query.Network, _, err = db.LookupNetwork(ip, &query); err != nil {
		s.logger.WithError(err).WithField("ip", ip.String()).Error("error looking up ip asn info")
		return nil, err
	}

	s.asnCache.Set(ip.String(), query)
	return query, nil
}

func (s *Service) lookupGeo(ctx context.Context, ip net.IP) (query *models.GeoQuery, err error) {
	metrics.LookupCount.WithLabelValues("geo").Inc()

	if val := s.geoCache.Get(ip.String()); val != nil {
		return val, nil
	}

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

	s.geoCache.Set(ip.String(), query)
	return query, nil
}
