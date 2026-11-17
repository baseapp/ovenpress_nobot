package utils

import "testing"

func TestStripRawQueryPrefix(t *testing.T) {
	cases := []struct{ in, want string }{
		{"", ""},
		{"wp-admin", "wp-admin"},
		{"wp-admin=", "wp-admin="},
		{"z=1&a=2", "z=1&a=2"},
		{"wp-admin&__goaway_token=abc", "wp-admin"},
		{"__goaway_referer=x&b=2&a&__goaway_id=1", "b=2&a"},
		{"q=a%20b&x", "q=a%20b&x"},
	}
	for _, c := range cases {
		if got := StripRawQueryPrefix(c.in, "__goaway_"); got != c.want {
			t.Errorf("StripRawQueryPrefix(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
