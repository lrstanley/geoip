package apihandler

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/chix"
)

var reLanguage = regexp.MustCompile(`^[^a-zA-Z, ]+.*?$`)

func (h *handler) getLookup(w http.ResponseWriter, r *http.Request) {
	addr := strings.TrimSpace(chi.URLParam(r, "addr"))
	filters := strings.Split(chi.URLParam(r, "filters"), ",")

	logger := chix.Log(r).WithField("lookup_addr", addr)

	// Prioritize "lang" query param.
	lang := h.lookupSvc.MatchLanguage(strings.ReplaceAll(r.FormValue("lang"), " ", ""))

	if lang == "" {
		// Try to get the language from the Accept-Language header.
		for _, l := range strings.Split(reLanguage.ReplaceAllString(r.Header.Get("Accept-Language"), ""), ",") {
			lang = h.lookupSvc.MatchLanguage(l)
			if lang != "" {
				break
			}
		}
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

	result, err := h.lookupSvc.IP(r.Context(), addr, filters, lang)
	if err != nil {
		logger.WithError(err).Error("error looking up addr")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	apiResponse(w, r, result, filters)
}
