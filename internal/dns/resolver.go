package dns

import (
	"context"
	"math/rand"
	"net"
	"strings"

	"github.com/lrstanley/geoip/internal/models"
)

var Resolver = net.DefaultResolver

func UseCustom(config models.ConfigDNS) {
	Resolver = &net.Resolver{PreferGo: true, Dial: newCustomResolver(config)}
}

type customResolver func(ctx context.Context, network, address string) (net.Conn, error)

func newCustomResolver(config models.ConfigDNS) customResolver {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		var index int

		if config.Local {
			index = rand.Intn(len(config.Resolvers) + 1)
		} else {
			// Generate a random number, which is used to select a resolver.
			// However, if the number generated is out of the bounds of the
			// amount of resolvers, use the system resolver, since they
			// requested it.
			index = rand.Intn(len(config.Resolvers))
		}

		if index == len(config.Resolvers) {
			return net.Dial(network, address)
		}

		addr := config.Resolvers[index]

		if strings.Contains(addr, ":") {
			return net.Dial(network, addr)
		}
		return net.Dial(network, addr+":53")
	}
}
