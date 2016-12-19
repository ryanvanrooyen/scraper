package scraper

import "testing"

func TestNew_GoodUrl(t *testing.T) {

	url := "url"
	getter := MemoryGetter{
		url: "<div>TestData</div>",
	}
	scraper := New(url, nil, getter)
	if scraper == nil {
		t.Fatalf("New created a nil scraper")
	}

	results, err := scraper.Done()

	if len(results) != 0 {
		t.Fatalf("Should be 0 results")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestNew_BadUrl(t *testing.T) {

	getter := MemoryGetter{}
	scraper := New("url", nil, getter)
	if scraper == nil {
		t.Fatalf("New created a nil scraper")
	}

	results, err := scraper.Done()
	if results != nil {
		t.Fatalf("Results should be nil")
	}
	if err == nil {
		t.Fatal(err)
	}
}
