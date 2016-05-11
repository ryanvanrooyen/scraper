package scraper

import (
	"io"

	"fmt"
	"strings"

	css "github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

// Scraper defines a simple
// web scraper's functionality
type Scraper interface {
	Filter(selector string) Scraper
	Select(selector Sel) Scraper
	Follow(selector string) Scraper
	Done() ([]map[string]string, error)
}

// Logger defines a simple interface for
// logging (compatible with std Log package)
type Logger interface {
	Printf(format string, v ...interface{})
}

// Sel (Selector) is a simple key-value map of
// prop namesto values based on a css selector
type Sel map[string]string

type scraper struct {
	Getter
	Logger  Logger
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
func New(url string, logger Logger, getter Getter) Scraper {

	if getter == nil {
		getter = HTTPGetter()
	}

	s := &scraper{Logger: logger, Getter: getter}

	resp, err := s.Get(url)
	if err != nil {
		return setError(s, err)
	}

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

func (s *scraper) Follow(selector string) Scraper {

	if s.Error != nil {
		return s
	}

	sel, err := css.Compile(selector)
	if err != nil {
		return setError(s, err)
	}

	var allNodes []*html.Node

	for _, node := range s.Nodes {
		urlNodes := sel.MatchAll(node)
		for _, urlNode := range urlNodes {
			if url, err := textOrAttr(selector, urlNode); err == nil {

				newNode, err := followURL(s, url)
				if err != nil {
					return setError(s, err)
				}
				allNodes = append(allNodes, newNode)
			}
		}
	}

	s.Nodes = allNodes
	return s
}

func followURL(s *scraper, url string) (*html.Node, error) {

	s.Logf("Getting url %s", url)
	rc, err := s.Get(url)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	result, err := html.Parse(rc)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *scraper) Done() ([]map[string]string, error) {

	if s.Error != nil {
		return nil, s.Error
	}

	return s.Results, nil
}

func (s *scraper) Logf(format string, a ...interface{}) {
	if s.Logger != nil {
		s.Logger.Printf(format+"\n", a)
	}
}

func selectText(selector string, node *html.Node) (string, error) {

	sel, err := css.Compile(selector)
	if err != nil {
		return "", err
	}

	result, nodes := "", sel.MatchAll(node)
	for _, n := range nodes {
		if txt, err := textOrAttr(selector, n); err == nil {
			result += txt
		} else {
			return "", err
		}
	}

	return result, nil
}

func get(s *scraper, r io.ReadCloser) Scraper {

	defer r.Close()
	result, err := html.Parse(r)
	if err != nil {
		return setError(s, err)
	}

	s.Nodes = make([]*html.Node, 1)
	s.Nodes[0] = result
	return s
}

func textOrAttr(selector string, node *html.Node) (string, error) {

	attrName := getAttrName(selector)
	if attrName == "" {
		return text(node)
	}

	return attr(node, attrName)
}

func getAttrName(selector string) string {

	startIndex := strings.LastIndex(selector, "[")
	endIndex := strings.LastIndex(selector, "]")

	if startIndex < 0 || endIndex < 0 ||
		startIndex >= endIndex ||
		endIndex != len(selector)-1 {
		return ""
	}

	attr := selector[startIndex+1 : endIndex]

	if !strings.Contains(attr, "=") {
		return attr
	}

	return attr[:strings.Index(attr, "=")]
}

func attr(node *html.Node, attr string) (string, error) {

	if len(node.Attr) == 0 {
		return "", fmt.Errorf(
			"No attributes found on node %v", node)
	}

	for _, a := range node.Attr {
		if strings.EqualFold(a.Key, attr) {
			return a.Val, nil
		}
	}

	return "", fmt.Errorf("No attribute %q found", attr)
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
