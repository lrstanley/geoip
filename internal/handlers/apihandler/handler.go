package apihandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lrstanley/chix"
	"github.com/lrstanley/geoip/internal/lookup"
)

type handler struct {
	lookupSvc *lookup.Service
}

func New(lookupSvc *lookup.Service) *handler {
	return &handler{
		lookupSvc: lookupSvc,
	}
}

func (h *handler) Route(r chi.Router) {
	r.Get("/{addr}", h.getLookup)
	r.Get("/{addr}/{filters}", h.getLookup)
	r.With(middleware.NoCache, middleware.GetHead).Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		chix.JSON(w, r, http.StatusOK, chix.M{"pong": true})
	})
}
