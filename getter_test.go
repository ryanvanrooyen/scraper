package scraper

import ("io/ioutil"; "testing")

func TestMultiUserAgent(t *testing.T) {

	agentStrings := []string{
		"userAgent1", "userAgent2", "userAgent3",
	}

	agent := &multiUserAgent{
		values: agentStrings,
	}

	for i := 0; i < len(agentStrings)*2; i++ {

		expected := agentStrings[i%len(agentStrings)]
		actual := agent.UserAgent()
		if actual != expected {
			t.Errorf("Expected %s, received %s", expected, actual)
		}
	}
}

func TestMemoryGetter(t *testing.T) {

	url := "localhost"
	expected := "test data"

	client := MemoryGetter{
		url: expected,
	}

	verifyGetter(t, client, url, expected)
}

func TestFileGetter(t *testing.T) {

	url := "localhost"
	fileName := "./testFiles/data.txt"
	fileData, err := ioutil.ReadFile(fileName)

	if err != nil {
		t.Error(err)
	}

	expected := string(fileData)
	client := FileGetter{
		url: fileName,
	}

	verifyGetter(t, client, url, expected)
}

func verifyGetter(t *testing.T, client Getter,
	url string, expectedData string) {

	reader, _ := client.Get(url)
	data, _ := ioutil.ReadAll(reader)
	actual := string(data)

	if actual != expectedData {
		t.Errorf("Expected \"%s\" received \"%s\"",
			expectedData, actual)
	}

	t.Logf("Expected matched actual %s", actual)
}
