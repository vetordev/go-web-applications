package main

import "testing"

func TestGetFilename(t *testing.T) {
	title := "GoWiki"
	filename := getFilename(title)

	want := title + ".txt"

	if filename != want {
		t.Fatalf(`getFilename(%q) = %q, want match for %q`, title, filename, want)
	}
}

func TestGetTitle(t *testing.T) {
	path := "/view/GoWiki"
	title, err := getTitle(path)

	want := "GoWiki"

	if title != want || err != nil {
		t.Fatalf(`getTitle(%q) = %q, %v, want match for %q, nil`, path, title, err, want)
	}
}

func TestGetTitleWithInvalidPath(t *testing.T) {
	path := "/view/$$%%##GoWiki"
	title, err := getTitle(path)

	want := ""

	if title != want || err == nil {
		t.Fatalf(`getTitle(%q) = %q, %q, want match for %q, "invalid page title"`, path, title, err, want)
	}
}
