// Package locale is reusable across Kurmon projects: the language set,
// URL prefixes (/en/, /ja/ ...) and canonical-path resolution.
// Behaviour matches kurmon-web/site/lang.go.
package locale

import (
	"path"
	"strings"
)

type Language struct{ Code, Segment, Name string }

// Languages is the shared language set.
// TODO(confirm): sixth language.
var Languages = []Language{
	{"en", "en", "English"},
	{"ja", "ja", "日本語"},
	{"ko", "ko", "한국어"},
	{"zh-Hant", "zh-hant", "繁體中文"},
	{"zh-Hans", "zh-hans", "简体中文"},
}

const Default = "en"

// Option is one entry of the language switcher (components/header/switch-lang).
type Option struct {
	Code, Name, URL string
	Current         bool
}

func FromSegment(segment string) (string, bool) {
	for _, l := range Languages {
		if l.Segment == strings.ToLower(segment) {
			return l.Code, true
		}
	}
	return "", false
}

func segment(code string) string {
	for _, l := range Languages {
		if l.Code == code {
			return l.Segment
		}
	}
	return Default
}

// Path("ja", "/checkin/") == "/ja/checkin/"
func Path(code, pagePath string) string {
	clean := strings.Trim(path.Clean("/"+pagePath), "/")
	if clean == "" || clean == "." {
		return "/" + segment(code) + "/"
	}
	return "/" + segment(code) + "/" + clean + "/"
}

func Options(code, pagePath string) []Option {
	out := make([]Option, 0, len(Languages))
	for _, l := range Languages {
		out = append(out, Option{Code: l.Code, Name: l.Name, URL: Path(l.Code, pagePath), Current: l.Code == code})
	}
	return out
}

func Current(code, pagePath string) Option {
	for _, o := range Options(code, pagePath) {
		if o.Current {
			return o
		}
	}
	return Options(Default, pagePath)[0]
}

// Resolve splits "/ja/checkin/" into ("ja", "/checkin/"). A non-empty
// redirect is the canonical URL the request should be sent to.
func Resolve(urlPath string) (code, pagePath, redirect string) {
	trimmed := strings.Trim(urlPath, "/")
	if trimmed == "" {
		return Default, "/", Path(Default, "/")
	}
	parts := strings.Split(trimmed, "/")
	if c, ok := FromSegment(parts[0]); ok {
		pagePath = "/"
		if len(parts) > 1 {
			pagePath = "/" + strings.Join(parts[1:], "/") + "/"
		}
		if canonical := Path(c, pagePath); canonical != urlPath {
			return c, pagePath, canonical
		}
		return c, pagePath, ""
	}
	pagePath = "/" + trimmed + "/"
	return Default, pagePath, Path(Default, pagePath)
}
