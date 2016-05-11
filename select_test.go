package scraper

import (
	"os"
	"testing"
	"log"
)

const selectTestHTML = `
	<div class="d1">
		<span class="a1"><span class="b1">V1</span></span>
		<span class="a1"><span class="b1 c">V2</span></span>
		<span class="a1"><span class="b1">V3</span></span>
	</div>
	<div class="d2">
		<span class="a1"><span class="b1">V1</span></span>
		<span class="a1"><span class="b1 c">V2</span></span>
		<span class="a1"><span class="b1">V3</span></span>
	</div>`

type SelectTest struct {
	Sel string
	Exp string
}

var selectTestData = []SelectTest{
	SelectTest{Sel: "span.a1",		Exp: "V1V2V3V1V2V3"},
	SelectTest{Sel: ".b1",			Exp: "V1V2V3V1V2V3"},
	SelectTest{Sel: ".d1 > span",	Exp: "V1V2V3"},
	SelectTest{Sel: "div",			Exp: "V1V2V3V1V2V3"},
	SelectTest{Sel: "div .c",		Exp: "V2V2"},
	SelectTest{Sel: ".c",			Exp: "V2V2"},
}

func TestSelect(t *testing.T) {
	for _, test := range selectTestData {
		selectTest(t, test.Sel, test.Exp)
	}
}

func selectTest(t *testing.T, selector string, exp string) {

	url := "localhost"
	getter := MemoryGetter{url: selectTestHTML}
	logger := log.New(os.Stdout, "", log.Lshortfile)
	scraper := New(url, logger, getter)

	if scraper == nil {
		t.Fatalf("New created a nil scraper")
	}

	propName := "value"
	scraper.Select(Sel{propName: selector})

	results, err := scraper.Done()
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, received %v: %v",
			len(results), results)
	}

	if act := results[0][propName]; exp != act {
		t.Fatalf("Expected %q, received %q: %v", exp, act, results)
	}
}
