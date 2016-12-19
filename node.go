package scraper

import (
	"fmt"
	"io"
	"strings"

	css "github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

type nodeFactory interface {
	Create(url string, r io.Reader) (node, error)
}

type node interface {
	Filter(sel string, cssSel css.Selector) []node
	Follow(sel string, cssSel css.Selector) []node
	Select(name string, sel string, cssSel css.Selector) error
	GetData() []map[string]string
}

type result struct {
	Getter
	URL     string
	Element *html.Node
	Data    []map[string]string
	Nodes   []*result
}

type nFactory struct {
	Getter
}

func (n nFactory) Create(url string, r io.Reader) (node, error) {

	el, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return &result{n.Getter, url, el, nil, nil}, nil
}

func (r *result) Filter(sel string, cssSel css.Selector) []node {

	if cssSel == nil {
		return []node{}
	}

	var nodes []node
	elements := cssSel.MatchAll(r.Element)

	for _, el := range elements {

		node := &result{
			r.Getter, r.URL, el, r.Data, nil,
		}

		r.Nodes = append(r.Nodes, node)
		nodes = append(nodes, node)
	}

	return nodes
}

func (r *result) Select(name string, sel string, cssSel css.Selector) error {

	nodes := cssSel.MatchAll(r.Element)
	r.Data = make([]map[string]string, len(nodes))

	for i, n := range nodes {

		txt, err := textOrAttr(sel, n)
		if err != nil {
			return err
		}

		dataMap := make(map[string]string)
		dataMap[name] = txt

		r.Data[i] = dataMap
	}

	return nil
}

func (r *result) Follow(sel string, cssSel css.Selector) []node {

	if cssSel == nil {
		return []node{}
	}

	urlNodes := cssSel.MatchAll(r.Element)
	nodes := make([]node, 0, len(urlNodes))

	for _, urlNode := range urlNodes {

		url, err := textOrAttr(sel, urlNode)
		if err != nil {
			continue
		}

		el, err := r.followURL(url)
		if err != nil {
			continue
		}
		nodes = append(nodes, &result{
			r.Getter, url, el, r.Data, nil,
		})
	}

	return nodes
}

func (r *result) GetData() []map[string]string {

	if r.Nodes == nil {
		return r.Data
	}

	var allData []map[string]string
	allData = append(allData, r.Data...)

	for _, n := range r.Nodes {
		allData = append(allData, n.GetData()...)
	}

	return allData
}

func selectText(selStr string, sel css.Selector, el *html.Node) (string, error) {

	txt, nodes := "", sel.MatchAll(el)
	for _, n := range nodes {

		t, err := textOrAttr(selStr, n)
		if err != nil {
			return "", err
		}
		txt += t
	}

	return txt, nil
}

func (r *result) followURL(url string) (*html.Node, error) {

	rc, err := r.Get(url, r.URL)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	el, err := html.Parse(rc)
	if err != nil {
		return nil, err
	}

	return el, nil
}

func textOrAttr(sel string, node *html.Node) (string, error) {

	attrName := getAttrName(sel)
	if attrName == "" {
		return text(node)
	}

	return attr(node, attrName)
}

func getAttrName(sel string) string {

	startIndex := strings.LastIndex(sel, "[")
	endIndex := strings.LastIndex(sel, "]")

	if startIndex < 0 || endIndex < 0 ||
		startIndex >= endIndex ||
		endIndex != len(sel)-1 {
		return ""
	}

	attr := sel[startIndex+1 : endIndex]

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
		txt = txt + strings.TrimSpace(t)
	}

	// Make sure to remove all white space from the txt
	return txt, nil
}
