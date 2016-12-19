package scraper

import (
	"io"

	css "github.com/andybalholm/cascadia"
)

// Logger defines a simple interface for
// logging (compatible with std Log package)
type Logger interface {
	Printf(format string, v ...interface{})
}

type nilLogger struct{}

type scraperLog struct {
	Logger
	*scraper
}

type getterLog struct {
	Logger
	Getter
}

type nodeLog struct {
	Logger
	node
}

type nFactoryLog struct {
	Logger
	nodeFactory
}

func (l nilLogger) Printf(format string, v ...interface{}) {
}

func (s *scraperLog) Filter(selector string) Scraper {
	s.Printf("Filtering %v nodes by %s\n", len(s.Nodes), selector)
	s.scraper.Filter(selector)
	return s
}

func (s *scraperLog) Done() ([]map[string]string, error) {
	s.Printf("Getting results\v")
	results, err := s.scraper.Done()
	if err != nil {
		s.Printf("Error getting results: %q\n", err)
	}
	s.Printf("Found a total of %v results\n", len(results))
	return results, err
}

func (g *getterLog) Get(url string, srcURL string) (io.ReadCloser, error) {
	g.Printf("Getting url %q with src url %q\n", url, srcURL)
	r, err := g.Getter.Get(url, srcURL)
	if err != nil {
		g.Printf("Error getting url %q: %q\n", url, err)
	}
	return r, err
}

func (n nFactoryLog) Create(url string, r io.Reader) (node, error) {
	newNode, err := n.nodeFactory.Create(url, r)
	if err != nil {
		n.Printf("Error creating node: %q\n", err)
		return nil, err
	}
	n.Printf("Created node for url %q\n", url)
	return nodeLog{n.Logger, newNode}, nil
}

func (n nodeLog) Filter(sel string, cssSel css.Selector) []node {
	nodes := n.node.Filter(sel, cssSel)
	n.Printf("Found %v nodes filtering by %s\n", len(nodes), sel)

	var results []node
	for _, node := range nodes {
		results = append(results, nodeLog{n.Logger, node})
	}
	return results
}

func (n nodeLog) Select(name string, sel string, cssSel css.Selector) error {

	n.Printf("Selecting data with selector %s", sel)
	return n.node.Select(name, sel, cssSel)
}
