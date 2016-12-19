package scraper

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Getter is HTTP get abstraction to enable
// a Scraper's data source to be changed
type Getter interface {
	Get(url string, srcURL string) (io.ReadCloser, error)
}

type userAgent interface {
	UserAgent() string
}

type httpGetter struct {
	userAgent
}

type multiUserAgent struct {
	values []string
	index  int
}

type urlResolver struct {
	Getter
}

// MemoryGetter is a Getter that
// retrieves strings by specified url key.
type MemoryGetter map[string]string

// FileGetter is a Getter that
// retrieves file data by specified url key.
type FileGetter map[string]string

// HTTPGetter create a Getter
// that retrieves urls over http.
func HTTPGetter() Getter {
	return urlResolver{httpGetter{
		&multiUserAgent{
			values: []string{
				"",
			},
		},
	}}
}

// Get looks up the string data for the specifed url
func (c MemoryGetter) Get(url string, srcURL string) (io.ReadCloser, error) {

	data := c[url]
	if data == "" {
		return nil, errors.New("No data provided for url: " + url)
	}

	return ioutil.NopCloser(strings.NewReader(data)), nil
}

// Get looks up the file data for the specifed url
func (c FileGetter) Get(url string, srcURL string) (io.ReadCloser, error) {

	fileName := c[url]
	if fileName == "" {
		return nil, errors.New("No file provided for url: " + url)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (u urlResolver) Get(urlStr string, srcURL string) (io.ReadCloser, error) {

	if srcURL == "" {
		return u.Getter.Get(urlStr, srcURL)
	}

	uri, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	base, err := url.Parse(srcURL)
	if err != nil {
		return nil, err
	}

	resolvedURI := base.ResolveReference(uri)
	return u.Getter.Get(resolvedURI.String(), "")
}

func (c httpGetter) Get(url string, srcURL string) (io.ReadCloser, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("UserAgent", c.UserAgent())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (source *multiUserAgent) UserAgent() string {
	value := source.values[source.index]
	source.index = (source.index + 1) % len(source.values)
	return value
}
