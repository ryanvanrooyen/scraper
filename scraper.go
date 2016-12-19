package scraper

import css "github.com/andybalholm/cascadia"

// Sel (Selector) is a simple key-value map of
// prop names to values based on a css selector
type Sel map[string]string

// Scraper defines a simple
// web scraper's functionality
type Scraper interface {
	Filter(selector string) Scraper
	Select(selector Sel) Scraper
	Follow(selector string) Scraper
	Done() ([]map[string]string, error)
}

type initializer interface {
	Scraper
	init(url string) Scraper
}

type scraper struct {
	Getter
	nodeFactory
	Nodes    []node
	RootNode node
	Error    error
}

// Get creates a new scraper by
// retrieving the HTML at the given URL
func Get(url string) Scraper {
	return New(url, nil, nil)
}

// New creates a new scraper by using
// the data provided by the specified Getter
func New(url string, logger Logger, getter Getter) Scraper {

	if getter == nil {
		getter = HTTPGetter()
	}
	var s initializer
	if logger != nil {
		getter = &getterLog{logger, getter}
		s = &scraperLog{
			logger,
			&scraper{
				getter,
				nFactoryLog{logger, nFactory{getter}},
				nil, nil, nil,
			},
		}
	} else {
		s = &scraper{
			getter,
			nFactory{getter},
			nil, nil, nil,
		}
	}

	return s.init(url)
}

func (s *scraper) Filter(selector string) Scraper {

	if s.Error != nil {
		return s
	}

	sel, err := css.Compile(selector)
	if err != nil {
		return s.setError(err)
	}

	var allNodes []node
	for _, n := range s.Nodes {
		nodes := n.Filter(selector, sel)
		allNodes = append(allNodes, nodes...)
	}

	s.Nodes = allNodes
	return s
}

func (s *scraper) Select(selectors Sel) Scraper {

	if s.Error != nil {
		return s
	}
	if selectors == nil {
		selectors = Sel{}
	}

	for name, selector := range selectors {
		sel, err := css.Compile(selector)
		if err != nil {
			return s.setError(err)
		}

		for _, n := range s.Nodes {
			err = n.Select(name, selector, sel)
			if err != nil {
				return s.setError(err)
			}
		}
	}

	return s
}

func (s *scraper) Follow(selector string) Scraper {

	if s.Error != nil {
		return s
	}

	sel, err := css.Compile(selector)
	if err != nil {
		return s.setError(err)
	}

	var allNodes []node
	for _, n := range s.Nodes {
		nodes := n.Follow(selector, sel)
		allNodes = append(allNodes, nodes...)
	}

	s.Nodes = allNodes
	return s
}

func (s *scraper) Done() ([]map[string]string, error) {

	if s.Error != nil {
		return nil, s.Error
	}
	return s.RootNode.GetData(), nil
}

func (s *scraper) init(url string) Scraper {

	resp, err := s.Get(url, "")
	if err != nil {
		return s.setError(err)
	}

	defer resp.Close()
	root, err := s.Create(url, resp)
	if err != nil {
		return s.setError(err)
	}

	s.RootNode = root
	s.Nodes = []node{root}
	return s
}

func (s *scraper) setError(err error) Scraper {
	s.Error = err
	return s
}
