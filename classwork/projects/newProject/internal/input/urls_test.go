package input

import "testing"

func TestNormalizateURL_AddsHTTPSwhenMissingSccheme(t *testing.T) {
	got:=NormalizeURL("example.com")
	want := "https://example.com"
	if got !=want{
		t.Fatalf("got %q, want %q",got,want)
	}
}

func TestNormalizateURL_LeavsExistingScheme(t *testing.T){
	got:=NormalizeURL("https://example.com")
	want := "https://example.com"
	if got !=want{
		t.Fatalf("got %q, want %q",got,want)
	}
}