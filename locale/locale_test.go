package locale

import "testing"

func TestResolve(t *testing.T) {
	cases := []struct{ in, code, page, redirect string }{
		{"/", "en", "/", "/en/"},
		{"/ja/", "ja", "/", ""},
		{"/zh-hans/checkin/", "zh-Hans", "/checkin/", ""},
		{"/ZH-HANS/checkin", "zh-Hans", "/checkin/", "/zh-hans/checkin/"},
		{"/checkin", "en", "/checkin/", "/en/checkin/"},
	}
	for _, c := range cases {
		code, page, redirect := Resolve(c.in)
		if code != c.code || page != c.page || redirect != c.redirect {
			t.Errorf("Resolve(%q) = %q %q %q", c.in, code, page, redirect)
		}
	}
}
