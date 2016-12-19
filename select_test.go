package scraper

import (
	"testing"
)

const selectTestHTML = `
	<div class="d1">
		<span class="a1"><span class="b1">V1</span></span>
		<span class="a2"><span class="b1 c">V2</span></span>
		<span class="a3"><span class="b1">V3</span></span>
	</div>
	<div class="d2">
		<span class="a1"><span class="b1">V1</span></span>
		<span class="a2"><span class="b1 c">V2</span></span>
		<span class="a3"><span class="b1">V3</span></span>
	</div>`

type SelectTest struct {
	Sel string
	Exp string
}

var selectTestData = []SelectTest{
	SelectTest{Sel: "span.a1", Exp: "V1V1"},
	SelectTest{Sel: ".b1", Exp: "V1V2V3V1V2V3"},
	SelectTest{Sel: ".d1 > span", Exp: "V1V2V3"},
	SelectTest{Sel: "div", Exp: "V1V2V3V1V2V3"},
	SelectTest{Sel: "div .c", Exp: "V2V2"},
	SelectTest{Sel: ".c", Exp: "V2V2"},
}

func testSelect(t *testing.T) {
	for _, test := range selectTestData {
		selectTest(t, test.Sel, test.Exp)
	}
}

func testMultiSelect1(t *testing.T) {

	scraper := createScraper()
	if scraper == nil {
		t.Fatalf("New created a nil scraper")
	}

	scraper.Select(Sel{
		"Prop1": ".d1 .a1",
		"Prop2": ".d1 .a2",
		"Prop3": ".d1 .a3",
	})

	verifySelectResults(t, scraper, []map[string]string{
		map[string]string{
			"Prop1": "V1",
			"Prop2": "V2",
			"Prop3": "V3",
		},
	})
}

func testMultiSelect2(t *testing.T) {

	scraper := createScraper()
	if scraper == nil {
		t.Fatalf("New created a nil scraper")
	}

	scraper.Select(Sel{"Prop1": ".d1 .a1"})
	scraper.Select(Sel{"Prop2": ".d1 .a2"})
	scraper.Select(Sel{"Prop3": ".d1 .a3"})

	verifySelectResults(t, scraper, []map[string]string{
		map[string]string{
			"Prop1": "V1",
			"Prop2": "V2",
			"Prop3": "V3",
		},
	})
}

func selectTest(t *testing.T, selector string, exp string) {

	scraper := createScraper()
	if scraper == nil {
		t.Fatalf("New created a nil scraper")
	}

	propName := "value"
	scraper.Select(Sel{propName: selector})

	verifySelectResults(t, scraper, []map[string]string{
		map[string]string{propName: exp},
	})
}

func createScraper() Scraper {

	url := "localhost"
	getter := MemoryGetter{url: selectTestHTML}
	var logger Logger //= log.New(os.Stdout, "", log.Lshortfile)
	return New(url, logger, getter)
}

func verifySelectResults(t *testing.T, s Scraper, results []map[string]string) {

	scrapedVals, err := s.Done()
	if err != nil {
		t.Fatal(err)
	}
	if len(scrapedVals) != len(results) {
		t.Fatalf("Expected %v result, received %v: %v",
			len(results), len(scrapedVals), scrapedVals)
	}

	for i, result := range results {
		for name, exp := range result {
			if act := scrapedVals[i][name]; exp != act {
				t.Fatalf("Expected %q, received %q: %v", exp, act, result)
			}
		}
	}
}
