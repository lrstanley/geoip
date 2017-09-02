package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

func registerAPI(r chi.Router) {
	r.Get("/api/{addr}", apiLookup)
	r.Get("/api/{addr}/{filter}", apiLookup)
}

func apiLookup(w http.ResponseWriter, r *http.Request) {
	addr := chi.URLParam(r, "addr")
	ip := net.ParseIP(addr)
	if ip == nil {
		ips, err := net.LookupHost(addr)
		if err != nil || len(ips) == 0 {
			debug.Printf("error looking up %q as host address: %s", addr, err)
			http.NotFound(w, r)
			return
		}

		ip = net.ParseIP(ips[0])
	}

	results, err := addrLookup(flags.DBPath, ip)
	if err != nil {
		debug.Printf("error looking up address %q (%q): %s", addr, ip, err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if filter := strings.Split(chi.URLParam(r, "filter"), ","); filter != nil && len(filter) > 0 && filter[0] != "" {
		base := make(map[string]*json.RawMessage)
		var tmp []byte

		tmp, err = json.Marshal(results)
		if err != nil {
			panic(err)
		}

		if err = json.Unmarshal(tmp, &base); err != nil {
			panic(err)
		}

		out := make([]string, len(filter))
		for i := 0; i < len(filter); i++ {
			out[i] = strings.Replace(fmt.Sprintf("%s", *base[filter[i]]), "\"", "", -1)
		}

		w.Write([]byte(strings.Join(out, "|")))
		return
	}

	enc := json.NewEncoder(w)

	if ok, _ := strconv.ParseBool(r.FormValue("pretty")); ok {
		enc.SetIndent("", "  ")
	}

	enc.SetEscapeHTML(false) // Otherwise the map url will get unicoded.
	err = enc.Encode(results)
	if err != nil {
		panic(err)
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
