package lookup

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/bluele/gcache"
	"github.com/lrstanley/geoip/internal/dns"
	"github.com/lrstanley/geoip/internal/models"
	"github.com/lrstanley/go-bogon"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

type Service struct {
	ctx    context.Context
	logger log.Interface
	config models.ConfigDB

	cache    gcache.Cache
	metadata models.Atomic[*maxminddb.Metadata]
}

func NewService(ctx context.Context, logger log.Interface, config models.ConfigDB) *Service {
	return &Service{
		ctx: ctx,
		logger: logger.WithFields(log.Fields{
			"src":  "lookup",
			"path": config.Path,
		}),
		config: config,
		cache:  gcache.New(config.CacheSize).ARC().Expiration(config.CacheExpire).Build(),
	}
}

// IP does a geoip lookup of an IP address. filters is passed into
// this function, in case there are any long running tasks which the user
// may not even want (e.g. reverse dns lookups).
func (s *Service) IP(ctx context.Context, addr string, filters []string, lang string) (*models.GeoResult, error) {
	var result *models.GeoResult
	var err error

	if lang == "" {
		lang = s.config.DefaultLanguage
	}

	cacheKey := addr + ":" + lang + ":" + strings.Join(filters, ":")
	if val, _ := s.cache.GetIFPresent(cacheKey); val != nil {
		result = val.(*models.GeoResult)
		result.Cached = true
		return result, nil
	}

	ip := net.ParseIP(addr)
	if ip == nil {
		var ips []string
		ips, err = net.LookupHost(addr)
		if err != nil || len(ips) == 0 {
			s.logger.WithError(err).Error("error looking up addr as hostname")

			return &models.GeoResult{Error: fmt.Sprintf("invalid ip/host specified: %s", addr)}, nil
		}

		ip = net.ParseIP(ips[0])
	}

	if is, _ := bogon.Is(ip.String()); is {
		return &models.GeoResult{Error: "internal address"}, nil
	}

	db, err := maxminddb.Open(s.config.Path)
	if err != nil {
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
		City:          query.City.Names[lang],
		Country:       query.Country.Names[lang],
		CountryCode:   query.Country.Code,
		Continent:     query.Continent.Names[lang],
		ContinentCode: query.Continent.Code,
		Lat:           query.Location.Lat,
		Long:          query.Location.Long,
		Timezone:      query.Location.TimeZone,
		PostalCode:    query.Postal.Code,
		Proxy:         query.Traits.Proxy,
	}

	var subdiv []string
	for i := 0; i < len(query.Subdivisions); i++ {
		subdiv = append(subdiv, query.Subdivisions[i].Names[lang])
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

	wantsHosts := len(filters) == 0
	if !wantsHosts {
		for i := 0; i < len(filters); i++ {
			if filters[i] == "host" {
				wantsHosts = true
				break
			}
		}
	}

	if wantsHosts {
		result.Host, _ = s.Reverse(ctx, ip)
	}

	if err = s.cache.Set(cacheKey, result); err != nil {
		s.logger.WithError(err).Error("error setting cache key for result")
	}

	return result, nil
}

func (s *Service) Reverse(ctx context.Context, addr net.IP) (string, error) {
	dnsCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var names []string
	var err error

	if names, err = dns.Resolver.LookupAddr(dnsCtx, addr.String()); err == nil && len(names) > 0 {
		return strings.TrimSuffix(names[0], "."), nil
	}

	return "", err
}

func (s *Service) UseMetadataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := s.metadata.Load()
		if m == nil {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("X-Maxmind-Build", fmt.Sprintf("%d-%d", m.IPVersion, m.BuildEpoch))
		w.Header().Set("X-Maxmind-Type", m.DatabaseType)

		next.ServeHTTP(w, r)
	})
}
