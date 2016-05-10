package scraper

import "testing"

func TestNew_GoodUrl(t *testing.T) {

	url := "localhost"
	getter := MemoryGetter{
		url: "<div>TestData</div>",
	}
	scraper := New(url, nil, getter)
	if scraper == nil {
		t.Fatalf("New created a nil scraper")
	}

	results, err := scraper.Done()
	if results != nil {
		t.Fatalf("Results should be nil")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestNew_BadUrl(t *testing.T) {

	getter := MemoryGetter{}
	scraper := New("localhost", nil, getter)
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