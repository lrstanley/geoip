// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package lookup

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/lrstanley/geoip/internal/models"
	maxminddb "github.com/oschwald/maxminddb-golang"
)

// NewUpdater returns a new service for monitoring database updates. If an update
// is needed, it will be downloaded, decompressed, verified, and installed.
func NewUpdater(config models.ConfigDB, logger log.Interface, lookupSvc *Service, dbType string) *Updater {
	updater := &Updater{
		config:    config,
		logger:    logger.WithField("src", fmt.Sprintf("updater-%s", dbType)),
		dbType:    dbType,
		lookupSvc: lookupSvc,
	}

	switch updater.dbType {
	case models.DatabaseGeoIP:
		updater.updateURL = fmt.Sprintf(config.GeoIPURL, config.LicenseKey)
		updater.path = config.GeoIPPath
	case models.DatabaseASN:
		updater.updateURL = fmt.Sprintf(config.ASNURL, config.LicenseKey)
		updater.path = config.ASNPath
	default:
		panic("unknown database type")
	}

	return updater
}

// Updater is a service for monitoring database updates.
type Updater struct {
	config    models.ConfigDB
	logger    log.Interface
	dbType    string
	updateURL string
	path      string

	lookupSvc *Service
}

// Start initiates checks for updates, and if an update is needed, it starts the
// update process.
func (u *Updater) Start(ctx context.Context) (err error) {
	var needsUpdate bool

	u.logger.Info("checking for database updates")
	needsUpdate, err = u.check()
	if err != nil {
		u.logger.WithError(err).Error("error checking current database status")
	}

	if needsUpdate {
		if err = u.update(); err != nil {
			u.logger.WithError(err).Error("unable to update database")
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(u.config.UpdateInterval):
			u.logger.Info("checking for database updates")

			needsUpdate, err = u.check()
			if err != nil {
				u.logger.WithError(err).Error("error checking current database status")
			}

			if needsUpdate {
				if err = u.update(); err != nil {
					u.logger.WithError(err).Error("unable to update database")
				}
			}
		}
	}
}

// updateMetadata updates the metadata information in the lookup service.
func (u *Updater) updateMetadata(path string) error {
	var err error
	var db *maxminddb.Reader
	var metadata maxminddb.Metadata

	db, err = maxminddb.Open(path)
	if err != nil {
		return err
	}
	defer db.Close()

	metadata = db.Metadata
	u.lookupSvc.Metadata.Store(u.dbType, &metadata)
	return nil
}

// check checks the current database status, and last modify date.
func (u *Updater) check() (needsUpdate bool, err error) {
	curSeconds := time.Now().UnixNano() / int64(time.Second)
	stat, err := os.Stat(u.path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}

		return true, err
	}

	err = u.updateMetadata(u.path)
	if err != nil {
		return false, err
	}

	if curSeconds-(stat.ModTime().UnixNano()/int64(time.Second)) < 604800 {
		return false, nil
	}

	return true, nil
}

// update downloads and verifies the database, then installs it.
func (u *Updater) update() error {
	started := time.Now()

	u.logger.Info("fetching new geoip data")

	archiveTempFile, err := ioutil.TempFile("", "geoip-archive-")
	if err != nil {
		return fmt.Errorf("unable to create temp file: %w", err)
	}

	defer func() {
		if err = archiveTempFile.Close(); err != nil {
			u.logger.WithError(err).WithField("fn", archiveTempFile.Name()).Error("error while closing file")
		}
		u.logger.WithField("fn", archiveTempFile.Name()).Info("deleting temp file")
		if err = os.Remove(archiveTempFile.Name()); err != nil {
			u.logger.WithError(err).WithField("fn", archiveTempFile.Name()).Error("error while removing temp file")
		}
	}()

	dbTempFile, err := ioutil.TempFile("", "geoip-db-")
	if err != nil {
		return fmt.Errorf("unable to create temp file: %w", err)
	}
	defer func() {
		if err = dbTempFile.Close(); err != nil {
			u.logger.WithError(err).WithField("fn", dbTempFile.Name()).Error("error while closing db temp file")
		}
		u.logger.WithField("fn", dbTempFile.Name()).Info("deleting db temp file")
		if err = os.Remove(dbTempFile.Name()); err != nil {
			u.logger.WithError(err).WithField("fn", dbTempFile.Name()).Error("error while removing db temp file")
		}
	}()

	u.logger.WithField("fn", dbTempFile.Name()).Info("streaming new database archive")
	resp, err := http.Get(u.updateURL)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gz)
	dbFound := false

	// loop through the tar file and extract first .mmdb file we find in the
	// archive.
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("unable to read tar archive: %w", err)
		}

		if header.Typeflag == tar.TypeReg && strings.HasSuffix(header.Name, ".mmdb") {
			u.logger.WithField("fn", dbTempFile.Name()).Info("found database in tar archive, extracting and writing")
			if _, err := io.Copy(dbTempFile, tarReader); err != nil {
				return fmt.Errorf("error extracting database from tar archive: %w", err)
			}
			dbFound = true
		}
	}

	if !dbFound {
		return errors.New("no database file found in tar archive")
	}

	if err := gz.Close(); err != nil {
		return err
	}

	u.logger.WithField("fn", dbTempFile.Name()).Info("successfully downloaded and decompressed new database, verifying")
	if _, err = dbTempFile.Seek(0, 0); err != nil {
		return err
	}

	db, err := maxminddb.Open(dbTempFile.Name())
	if err != nil {
		return fmt.Errorf("error while attempting to verify geoip data: %s", err)
	}

	if err = db.Verify(); err != nil {
		db.Close()
		return fmt.Errorf("error while attempting to verify geoip data: %s", err)
	}
	db.Close()

	err = u.updateMetadata(dbTempFile.Name())
	if err != nil {
		return err
	}

	u.logger.Info("verification complete, updating active database")

	file, err := os.Create(u.path)
	if err != nil {
		return err
	}

	var written int64
	written, err = io.Copy(file, dbTempFile)
	if err != nil {
		return err
	}

	u.logger.WithFields(log.Fields{
		"fn":       file.Name(),
		"bytes":    written,
		"duration": time.Since(started),
	}).Info("successfully wrote database update")

	return nil
}
