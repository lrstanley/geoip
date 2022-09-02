// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package httpware

import (
	"net/http"
	"regexp"
	"strings"
)

var reLanguage = regexp.MustCompile(`^[^a-zA-Z, ]+.*?$`)

func GetLanguage(r *http.Request) (languages []string) {
	// Prioritize "lang" query param.
	if lang := strings.ReplaceAll(r.FormValue("lang"), " ", ""); lang != "" {
		languages = append(languages, lang)
	}

	// Try to get the language from the Accept-Language header.
	for _, l := range strings.Split(reLanguage.ReplaceAllString(r.Header.Get("Accept-Language"), ""), ",") {
		if l != "" {
			languages = append(languages, l)
		}
	}

	return languages
}
