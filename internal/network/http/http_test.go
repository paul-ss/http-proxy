package http

import (
	"fmt"
	"testing"
)

var reqUrl = []struct {
	in string
	out string
}{
	{"google.com", ""},
	{"http://google.com", ""},
	{"https://google.com", ""},
	{"google.com:443", ""},
	{"google.com/jopa", ""},
	{"http://google.com/jopa", ""},
	{"https://google.com/jopa", ""},
	{"google.com:443/jopa", ""},
}

func TestFirstLineUrl(t *testing.T) {
	for _, tt := range reqUrl {
		t.Run(tt.in, func(t *testing.T) {
			tt.in = "1 " + tt.in + " 1"
			r := Request{}
			r.parseFirstLine([]byte(tt.in))
			fmt.Printf("Scheme: %s, Host: %s, Port: %s, Path: %s Url: %s\n",
				r.Url.Scheme,
				r.Url.Host,
				r.Url.Port(),
				r.Url.Path,
				r.Url.String())
		})
	}
}

