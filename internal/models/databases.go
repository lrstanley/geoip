// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

const (
	DatabaseGeoIP = "geoip"
	DatabaseASN   = "asn"
)

var Databases = []string{DatabaseGeoIP, DatabaseASN}
