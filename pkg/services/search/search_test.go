package search

import "testing"

func TestFormatURL(t *testing.T) {
	want := "http://www.google.com"
	got := formatURL("www.google.com")
	if want != got {
		t.Errorf("Wanted correctly formatted URL: %s, got %s", want, got)
	}
	want = "http://youtube.com"
	got = formatURL("youtube.com")
	if want != got {
		t.Errorf("Wanted correctly formatted URL: %s, got %s", want, got)
	}
}
