package scraper

import ("testing")

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
