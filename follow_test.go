package scraper

import ("testing"; "os"; "log")

const mainHTML = `
	<div class="d1">
		<span class="a1"><a href="/page1">Page1Link</a></span>
	</div>
	<div class="d2">
		<span class="a1"><a href="/page2">Page1Link</a></span>
	</div>`

const page1HTML = `
	<div class="p1">
		<span class="a1"><span class="b1">P1Data1</span></span>
		<span class="a2"><span class="b2">P1Data2</span></span>
		<span class="a3"><span class="b3">P1Data3</span></span>
	</div>`

const page2HTML = `
	<div class="p2">
		<span class="a1"><span class="b1">P2Data1</span></span>
		<span class="a2"><span class="b2">P2Data2</span></span>
		<span class="a3"><span class="b3">P2Data3</span></span>
	</div>`

type FollowTest struct {
	Fol string
	Sel string
	Exp string
}

var followTestData = []FollowTest{
	FollowTest{Fol: ".d1 a[href]", 		Sel: ".p1",	Exp: "P1Data1P1Data2P1Data3"},
	FollowTest{Fol: ".d1 a[href]", 		Sel: ".p1 .b2",	Exp: "P1Data2"},
	FollowTest{Fol: ".d2 span a[href]", Sel: ".p2 .a3",	Exp: "P2Data3"},
}

func TestFollow(t *testing.T) {
	for _, test := range followTestData {
		followTest(t, test.Fol, test.Sel, test.Exp)
	}
}

func followTest(t *testing.T, followSel string,
	selector string, expected string) {

	getter := MemoryGetter{
		"localhost":	mainHTML,
		"/page1":		page1HTML,
		"/page2":		page2HTML,
	}
	logger := log.New(os.Stdout, "", log.Lshortfile)
	scraper := New("localhost", logger, getter)

	if scraper == nil {
		t.Fatalf("New created a nil scraper")
	}

	scraper.Follow(followSel)

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

	if act := results[0][propName]; expected != act {
		t.Fatalf("Expected %s, received %s: %v", expected, act, results)
	}
}