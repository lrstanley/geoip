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
	maxminddb "github.com/oschwald/maxminddb-golang"
)

func (s *Service) Updater(ctx context.Context) error {
	var needsUpdate bool
	var err error

	s.logger.Info("checking for database updates")
	needsUpdate, err = s.checkForUpdates()
	if err != nil {
		s.logger.WithError(err).Error("error checking current database status")
	}

	if needsUpdate {
		if err = s.doUpdate(); err != nil {
			s.logger.WithError(err).Error("unable to update database")
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.config.UpdateInterval):
			s.logger.Info("checking for database updates")
			needsUpdate, err = s.checkForUpdates()
			if err != nil {
				s.logger.WithError(err).Error("error checking current database status")
			}

			if needsUpdate {
				if err = s.doUpdate(); err != nil {
					s.logger.WithError(err).Error("unable to update database")
				}
			}
		}
	}
}

func (s *Service) checkForUpdates() (needsUpdate bool, err error) {
	curSeconds := time.Now().UnixNano() / int64(time.Second)
	stat, err := os.Stat(s.config.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}

		return true, err
	}

	var db *maxminddb.Reader

	db, err = maxminddb.Open(s.config.Path)
	if err != nil {
		return true, err
	}

	s.metadata.Store(&db.Metadata)

	if curSeconds-(stat.ModTime().UnixNano()/int64(time.Second)) < 604800 {
		return false, nil
	}

	return true, nil
}

func (s *Service) doUpdate() error {
	started := time.Now()
	url := fmt.Sprintf(s.config.UpdateURL, s.config.LicenseKey)

	s.logger.Info("fetching new geoip data")

	archiveTempFile, err := ioutil.TempFile("", "geoip-archive-")
	if err != nil {
		return fmt.Errorf("unable to create temp file: %w", err)
	}

	defer func() {
		if err = archiveTempFile.Close(); err != nil {
			s.logger.WithError(err).WithField("fn", archiveTempFile.Name()).Error("error while closing file")
		}
		s.logger.WithField("fn", archiveTempFile.Name()).Info("deleting temp file")
		if err = os.Remove(archiveTempFile.Name()); err != nil {
			s.logger.WithError(err).WithField("fn", archiveTempFile.Name()).Error("error while removing temp file")
		}
	}()

	dbTempFile, err := ioutil.TempFile("", "geoip-db-")
	if err != nil {
		return fmt.Errorf("unable to create temp file: %w", err)
	}
	defer func() {
		if err = dbTempFile.Close(); err != nil {
			s.logger.WithError(err).WithField("fn", dbTempFile.Name()).Error("error while closing db temp file")
		}
		s.logger.WithField("fn", dbTempFile.Name()).Info("deleting db temp file")
		if err = os.Remove(dbTempFile.Name()); err != nil {
			s.logger.WithError(err).WithField("fn", dbTempFile.Name()).Error("error while removing db temp file")
		}
	}()

	s.logger.WithField("fn", dbTempFile.Name()).Info("streaming new database archive")
	resp, err := http.Get(url)
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
			s.logger.WithField("fn", dbTempFile.Name()).Info("found database in tar archive, extracting and writing")
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

	s.logger.WithField("fn", dbTempFile.Name()).Info("successfully downloaded and decompressed new database, verifying")
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

	s.metadata.Store(&db.Metadata)
	db.Close()

	s.logger.Info("verification complete, updating active database")

	file, err := os.Create(s.config.Path)
	if err != nil {
		return err
	}

	var written int64
	written, err = io.Copy(file, dbTempFile)
	if err != nil {
		return err
	}

	s.logger.WithFields(log.Fields{
		"fn":       file.Name(),
		"bytes":    written,
		"duration": time.Since(started),
	}).Info("successfully wrote database update")

	return nil
}
