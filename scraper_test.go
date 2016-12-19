package scraper

import (
	"testing"
	"fmt"
)

func ExampleGet() {

	results, err :=
		Get("https://golang.org/").
		Select(Sel{
			"title": "#heading-wide a",
		}).
		Follow(".read a[href]").
		Select(Sel{
			"blogTitle": "h1 a",
		}).
		Done()

	if err != nil {
		panic(err)
	}

	fmt.Println(len(results))
	fmt.Println(results[0]["title"])
	fmt.Println(results[0]["blogTitle"])
	/// Output:
	// 2
	// The Go Programming Language
	// The Go Blog
}

func TestAttrParsing(t *testing.T) {

	selectors := map[string]string {
		""						: "",
		"a"						: "",
		"a["					: "",
		"a]"					: "",
		"a[]"					: "",
		"a[href]"				: "href",
		"a[href] span"			: "",
		`a[href="/url"]`		: "href",
		"div[myAttr]"			: "myAttr",
		`div[myAttr="value"]`	: "myAttr",
		`div[myAttr="value"] p`	: "",
	}

	for selector, expected := range selectors {
		actual := getAttrName(selector)
		if expected != actual {
			t.Errorf("Parsing %q, actual %q expected %q",
				selector, actual, expected)
		}
	}
}
