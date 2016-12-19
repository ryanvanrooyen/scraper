package scraper

import (
	"io"
	"testing"
)

type URLTest struct {
	URL    string
	SrcURL string
	Exp    string
}

var urlTests = []URLTest{
	URLTest{URL: "http://test", SrcURL: "", Exp: "http://test"},
	URLTest{URL: "/p1", SrcURL: "", Exp: "/p1"},
	URLTest{URL: "/p1", SrcURL: "http://test", Exp: "http://test/p1"},
}

func TestUrlResolver(t *testing.T) {

	for _, test := range urlTests {
		getter := urlResolver{testGetter{test.Exp, t}}
		getter.Get(test.URL, test.SrcURL)
	}
}

type testGetter struct {
	expectedURL string
	t           *testing.T
}

func (g testGetter) Get(url string, srcURL string) (io.ReadCloser, error) {

	if g.expectedURL == url {
		return nil, nil
	}

	g.t.Errorf("Expected url %q received %q", g.expectedURL, url)
	return nil, nil
}
