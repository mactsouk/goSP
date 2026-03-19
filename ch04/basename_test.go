package main

import (
	"testing"
)

func TestBasename(t *testing.T) {
	tests := []struct {
		path     string
		suffixes []string
		expected string
	}{
		// Basic filenames
		{"file.txt", nil, "file.txt"},
		{"file.txt", []string{".txt"}, "file"},
		{"file.txt", []string{".log"}, "file.txt"}, // suffix not found

		// Paths with directories
		{"/usr/local/bin/go", nil, "go"},
		{"/usr/local/bin/go", []string{"go"}, ""}, // suffix strips full match
		{"/usr/local/bin/go1.20", []string{"1.20"}, "go"},

		// Trailing slashes
		{"/usr/local/bin/", nil, "bin"},
		{"////", nil, "/"}, // multiple slashes collapse to "/"

		// Root directory
		{"/", nil, "/"},

		// Hidden files
		{".gitignore", nil, ".gitignore"},
		{".gitignore", []string{".gitignore"}, ""}, // full suffix removal

		// Mixed cases
		{"myfile.go.txt", []string{".txt"}, "myfile.go"},
		{"myfile.go.txt", []string{".go.txt"}, "myfile"},
	}

	for _, tc := range tests {
		result := basename(tc.path, tc.suffixes...)
		if result != tc.expected {
			t.Errorf("basename(%q, %q) = %q; want %q",
				tc.path, tc.suffixes, result, tc.expected)
		}
	}
}
