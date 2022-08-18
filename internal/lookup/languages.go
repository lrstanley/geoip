// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package lookup

import "strings"

// var supportedLanguages = []string{
// 	"de", "en", "es", "fr", "ja", "pt-BR", "ru", "zh-CN",
// }

func (s *Service) MatchLanguage(lang string) (match string) {
	if lang == "" {
		return ""
	}

	supported := s.metadata.Load().Languages

	for i := 0; i < len(supported); i++ {
		if strings.EqualFold(lang, supported[i]) {
			return supported[i]
		}

		if j := strings.Index(supported[i], "-"); j > 0 {
			if strings.EqualFold(lang, supported[i][:j]) {
				return supported[i]
			}

			if k := strings.Index(lang, "-"); k > 0 {
				if strings.EqualFold(lang[:k], supported[i][:j]) {
					return supported[i]
				}
			}
		}

		if j := strings.Index(lang, "-"); j > 0 {
			if strings.EqualFold(lang[:j], supported[i]) {
				return supported[i]
			}
		}
	}

	return ""
}
