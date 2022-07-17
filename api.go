// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/bluele/gcache"
	"github.com/go-chi/chi"
	bogon "github.com/lrstanley/go-bogon"
)

func registerAPI(r chi.Router) {
	r.Get("/api/{addr}", apiLookup)
	r.Get("/api/{addr}/{filters}", apiLookup)
}

var reLanguage = regexp.MustCompile(`^[^a-zA-Z, ]+.*?$`)

func apiLookup(w http.ResponseWriter, r *http.Request) {
	addr := strings.TrimSpace(chi.URLParam(r, "addr"))
	filters := strings.Split(chi.URLParam(r, "filters"), ",")

	// Prioritize "lang" query param.
	lang := matchLanguage(strings.ReplaceAll(r.FormValue("lang"), " ", ""))

	if lang == "" {
		// Try to get the language from the Accept-Language header.
		for _, l := range strings.Split(reLanguage.ReplaceAllString(r.Header.Get("Accept-Language"), ""), ",") {
			lang = matchLanguage(l)
			if lang != "" {
				break
			}
		}
	}

	// If we still don't have a language, default to global language default.
	if lang == "" {
		lang = flags.DefaultLanguage
	}

	// If they're trying to send us way too many filters (which could cause
	// unwanted extra memory usage/be considered a resource usage attack),
	// we shouldn't handle their request.
	if len(filters) > 20 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: too many filters supplied")
		return
	}

	if len(filters) == 1 && filters[0] == "" {
		filters = []string{}
	}

	// Allow users to query themselves without having to have them specify
	// their own IP address. Note that this will not work if you are querying
	// the IP address locally.
	if self := strings.ToLower(addr); self == "self" || self == "me" {
		if strings.Contains(r.RemoteAddr, ":") {
			addr, _, _ = net.SplitHostPort(r.RemoteAddr)
		} else {
			addr = r.RemoteAddr
		}
	}

	// This would be the index key used for arc cache, if they request custom
	// filters, we should add that to the key, because those filters may
	// mean that the returned lookup has excluded information, which may
	// cause issues if the same query is returned with no requested filters.
	key := addr
	if len(filters) > 0 {
		key = addr + ":" + strings.Join(filters, ",")
	}

	key += "," + lang

	var result *AddrResult

	query, err := arc.GetIFPresent(key)
	if err == nil {
		resultFromARC, _ := query.(AddrResult)
		result = &resultFromARC
		w.Header().Set("X-Cache", "HIT")
		logger.Printf("query %s fetched from arc cache", addr)

		apiResponse(w, r, result, filters)
		return
	}

	w.Header().Set("X-Cache", "MISS")
	if err != gcache.KeyNotFoundError {
		logger.Printf("unable to get %s off arc stack: %s", addr, err)
	}

	ip := net.ParseIP(addr)
	if ip == nil {
		var ips []string
		ips, err = net.LookupHost(addr)
		if err != nil || len(ips) == 0 {
			logger.Printf("error looking up %q as host address: %s", addr, err)

			result = &AddrResult{Error: fmt.Sprintf("invalid ip/host specified: %s", addr)}
			apiResponse(w, r, result, filters)
			return
		}

		ip = net.ParseIP(ips[0])
	}

	if is, _ := bogon.Is(ip.String()); is {
		result = &AddrResult{Error: "internal address"}
		apiResponse(w, r, result, filters)
		return
	}

	result, err = addrLookup(r.Context(), ip, filters, lang)
	if err != nil {
		logger.Printf("error looking up address %q (%q): %s", addr, ip, err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if err = arc.Set(key, *result); err != nil {
		logger.Printf("unable to add %s to arc cache: %s", addr, err)
	}

	apiResponse(w, r, result, filters)
}

func apiResponse(w http.ResponseWriter, r *http.Request, result *AddrResult, filters []string) {
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

		out := make([]string, len(filters))
		for i := 0; i < len(filters); i++ {
			out[i] = strings.ReplaceAll(fmt.Sprintf("%s", *base[filters[i]]), "\"", "")
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strings.Join(out, "|")))
		return
	}

	enc := json.NewEncoder(w)

	if ok, _ := strconv.ParseBool(r.FormValue("pretty")); ok {
		enc.SetIndent("", "  ")
	}

	enc.SetEscapeHTML(false) // Otherwise the map url will get unicoded.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = enc.Encode(result)
	if err != nil {
		logger.Printf("error during json encode for %s: %s", r.RemoteAddr, err)
	}
}

func dbDetailsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mcache.RLock()
		if mcache.cache == nil {
			mcache.RUnlock()
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("X-Maxmind-Build", fmt.Sprintf("%d-%d", mcache.cache.IPVersion, mcache.cache.BuildEpoch))
		w.Header().Set("X-Maxmind-Type", mcache.cache.DatabaseType)
		mcache.RUnlock()

		next.ServeHTTP(w, r)
	})
}
