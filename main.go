// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"os"

	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/clix/v2"
	"github.com/lrstanley/geoip/internal/dns"
	"github.com/lrstanley/geoip/internal/lookup"
	"github.com/lrstanley/geoip/internal/models"
	"github.com/lrstanley/x/sync/scheduler"
)

// Should be automatically added by goreleaser.
const (
	version = "master"
	commit  = "master"
	date    = "unknown"
)

var (
	cli       *clix.CLI[models.Flags]
	lookupSvc *lookup.Service
)

func main() {
	cli = clix.NewWithDefaults(
		clix.WithAppInfo[models.Flags](clix.AppInfo{
			Name:        "geoip",
			Description: "geoip lookup service and api",
			Version:     version,
			Commit:      commit,
			Date:        date,
			Links:       clix.GithubLinks("github.com/lrstanley/geoip", "master", "https://liam.sh"),
		}),
	)

	logger := cli.GetLogger()
	ctx := context.Background()

	resolver := dns.NewResolver(ctx, cli.Flags.DNS)
	lookupSvc = lookup.NewService(ctx, logger, cli.Flags.DB, resolver)

	geoIPUpdater := lookup.NewUpdater(cli.Flags.DB, logger, lookupSvc, models.DatabaseGeoIP)
	asnUpdater := lookup.NewUpdater(cli.Flags.DB, logger, lookupSvc, models.DatabaseASN)

	if err := chix.Run(
		ctx,
		logger,
		httpServer(ctx, logger),
		scheduler.JobFunc(geoIPUpdater.Start),
		scheduler.JobFunc(asnUpdater.Start),
	); err != nil {
		logger.ErrorContext(ctx, "shutting down", "error", err)
		os.Exit(1)
	}
}
