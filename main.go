// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"

	"github.com/apex/log"
	"github.com/lrstanley/chix"
	"github.com/lrstanley/clix"
	"github.com/lrstanley/geoip/internal/dns"
	"github.com/lrstanley/geoip/internal/lookup"
	"github.com/lrstanley/geoip/internal/models"
)

// Should be automatically added by goreleaser.
var (
	version = "master"
	commit  = "master"
	date    = "unknown"
)

var (
	logger log.Interface
	cli    = &clix.CLI[models.Flags]{
		Links: clix.GithubLinks("github.com/lrstanley/geoip", "master", "https://liam.sh"),
	}

	lookupSvc *lookup.Service
)

func main() {
	cli.Parse()
	logger = cli.Logger

	if len(cli.Flags.DNS.Resolvers) > 0 {
		dns.UseCustom(cli.Flags.DNS)
	}

	ctx := log.NewContext(context.Background(), logger)

	lookupSvc = lookup.NewService(ctx, logger, cli.Flags.DB)

	if err := chix.RunCtx(
		ctx,
		httpServer(ctx),
		lookupSvc.Updater,
	); err != nil {
		logger.WithError(err).Fatal("shutting down")
	}
}
