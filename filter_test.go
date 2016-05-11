package scraper

import ("testing"; "os"; "log")

const filterTestHTML = `
	<div class="d1">
		<span class="a1"><span class="b1">V1</span></span>
		<span class="a1"><span class="b1">V2</span></span>
		<span class="a1"><span class="b1">V3</span></span>
	</div>
	<div class="d2">
		<span class="a1"><span class="b1">V1</span></span>
		<span class="a1"><span class="b1">V2</span></span>
		<span class="a1"><span class="b1">V3</span></span>
	</div>`

type FilterTest struct {
	Fil string
	Sel string
	Exp []string
}

var filterTestData = []FilterTest{
	FilterTest{Fil: ".d1 > span", Sel: ".b1",	Exp: []string{"V1","V2","V3"}},
}

func TestFilter(t *testing.T) {
	for _, test := range filterTestData {
		filterTest(t, test.Fil, test.Sel, test.Exp)
	}
}

func filterTest(t *testing.T, filter string,
	selector string, expected []string) {

	url := "localhost"
	getter := MemoryGetter{url:filterTestHTML}
	logger := log.New(os.Stdout, "", log.Lshortfile)
	scraper := New(url, logger, getter)

	if scraper == nil {
		t.Fatalf("New created a nil scraper")
	}

	scraper.Filter(filter)

	propName := "value"
	scraper.Select(Sel{propName: selector})

	results, err := scraper.Done()
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != len(expected) {
		t.Fatalf("Expected %v results, received %v: %v",
			len(expected), len(results), results)
	}

	for i := 0; i < len(results); i++ {
		exp, act := expected[i], results[i][propName]
		if exp != act {
			t.Fatalf("Expected %s, received %s: %v", exp, act, results)
		}
	}
}