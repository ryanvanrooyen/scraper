package scraper

import "testing"

const filterTestHTML = `
	<div class="d1">
		<span class="a1"><span class="b1">X1</span></span>
		<span class="a2"><span class="b1">X2</span></span>
		<span class="a3"><span class="b1">X3</span></span>
	</div>
	<div class="d2">
		<span class="a1"><span class="b1">Y1</span></span>
		<span class="a2"><span class="b1">Y2</span></span>
		<span class="a3"><span class="b1">Y3</span></span>
	</div>`

type FilterTest struct {
	Fil string
	Sel string
	Exp []string
}

var filterTests = []FilterTest{
	FilterTest{
		Fil: ".d1",
		Sel: "span.b1",
		Exp: []string{"X1", "X2", "X3"}},
	FilterTest{
		Fil: ".d1 > span",
		Sel: ".b1",
		Exp: []string{"X1", "X2", "X3"}},
	FilterTest{
		Fil: ".d2",
		Sel: ".a1",
		Exp: []string{"Y1"}},
}

func TestFilter(t *testing.T) {
	for _, test := range filterTests {
		filterTest(t, test.Fil, test.Sel, test.Exp)
	}
}

func filterTest(t *testing.T, filter string,
	selector string, expected []string) {

	url := "localhost"
	getter := MemoryGetter{url: filterTestHTML}
	// logger := log.New(os.Stdout, "TestFilter: ", log.LUTC)
	scraper := New(url, nil, getter)

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
