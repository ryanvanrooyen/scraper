package scraper

import (
	"fmt"
	"io"

	css "github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	"strings"
)

// Scraper defines a simple
// web scraper's functionality
type Scraper interface {
	Filter(selector string) Scraper
	Select(selector Sel) Scraper
	Done() ([]map[string]string, error)
}

// Sel (Selector) is a simple key-value map of
// prop namesto values based on a css selector
type Sel map[string]string

type scraper struct {
	Getter
	Logger  io.Writer
	Nodes   []*html.Node
	Results []map[string]string
	Error   error
}

// Get creates a new scraper by
// retrieving the HTML at the given URL
func Get(url string) Scraper {
	return New(url, nil, nil)
}

// New creates a new scraper by using
// the data provided by the given HTTPClient
func New(url string, logger io.Writer, getter Getter) Scraper {

	if getter == nil {
		getter = HTTPGetter()
	}

	s := &scraper{Logger: logger, Getter: getter}

	resp, err := s.Get(url)
	if err != nil {
		return setError(s, err)
	}

	defer resp.Close()
	return get(s, resp)
}

func (s *scraper) Filter(selector string) Scraper {

	if s.Error != nil {
		return s
	}

	sel, err := css.Compile(selector)
	if err != nil {
		return setError(s, err)
	}

	var allNodes []*html.Node

	for _, res := range s.Nodes {
		results := sel.MatchAll(res)
		for _, result := range results {
			if txt, err := text(result); err == nil {
				s.Logf("Found result %s", txt)
			}
		}
		allNodes = append(allNodes, results...)
	}

	s.Nodes = allNodes
	return s
}

func (s *scraper) Select(selectors Sel) Scraper {

	if s.Error != nil || selectors == nil {
		return s
	}

	var results []map[string]string

	for _, n := range s.Nodes {

		result := make(map[string]string)

		for prop, sel := range selectors {
			if val, err := selectText(sel, n); err == nil {
				result[prop] = val
			}
		}

		results = append(results, result)
	}

	s.Results = results
	return s
}

func (s *scraper) Done() ([]map[string]string, error) {

	if s.Error != nil {
		return nil, s.Error
	}

	return s.Results, nil
}

func (s *scraper) Logf(format string, a ...interface{}) {
	if s.Logger != nil {
		fmt.Fprintf(s.Logger, format+"\n", a)
	}
}

func selectText(selector string, node *html.Node) (string, error) {

	sel, err := css.Compile(selector)
	if err != nil {
		return "", err
	}

	result, nodes := "", sel.MatchAll(node)
	for _, n := range nodes {
		if txt, err := text(n); err == nil {
			result += txt
		} else {
			return "", err
		}
	}

	return result, nil
}

func get(s *scraper, r io.Reader) Scraper {

	result, err := html.Parse(r)
	if err != nil {
		return setError(s, err)
	}

	s.Nodes = make([]*html.Node, 1)
	s.Nodes[0] = result
	return s
}

func text(node *html.Node) (string, error) {

	if node.Type == html.TextNode {
		return node.Data, nil
	}

	txt := ""
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		t, _ := text(n)
		txt = txt + t
	}

	// Make sure to remove all white space from the txt
	return strings.Join(strings.Fields(txt), ""), nil
}

func setError(s *scraper, err error) Scraper {
	s.Error = err
	return s
}
