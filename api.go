package main

import (
	"net"
	"net/http"

	"github.com/go-chi/chi"
)

func registerAPI(r chi.Router) {
	r.Get("/api/{addr}", apiLookup)
	r.Get("/api/{addr}/{type}", apiLookup)
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

	// TODO: results.
	_, err := addrLookup(flags.DBPath, ip)
	if err != nil {
		debug.Printf("error looking up address %q (%q): %s", addr, ip, err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.Write([]byte("TODO\n"))
}
