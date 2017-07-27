package ep

import "testing"

func TestCantPlayUnsupportedFormat(t *testing.T) {
	fakeUrl := "fake/url/file.examp2225645le"
	err := Stream(Episode{URL:fakeUrl})
	if err == nil {
		t.Log("No error thrown for bad url:", fakeUrl)
		t.Fail()
	}
}

func TestGetExtension(t *testing.T) {
	fakeUrlWithLotsOfDots := `http://www.some.long.url.with.lots.of.dots.co.uk/mp3/audio.mp3`
	ext, err := getExtension(fakeUrlWithLotsOfDots)
	if err != nil {
		t.Log("Error in getExtension():", err.Error())
		t.Fail()
	}
	if ext != "mp3" {
		t.Log("Returned extension:", ext, "expected mp3")
	}
}