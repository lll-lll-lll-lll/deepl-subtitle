package vtt

import (
	"testing"
)

func TestCheckToken(t *testing.T) {
	// Test case for StartOrEndTime pattern
	if !checkToken(StartOrEndTime, "10:00") {
		t.Error("Expected true, got false")
	}

	// Test case for Separator pattern
	if !checkToken(Separator, "-->") {
		t.Error("Expected true, got false")
	}

	// Test case for Position pattern
	if !checkToken(Position, "position:50%") {
		t.Error("Expected true, got false")
	}

	// Test case for Line pattern
	if !checkToken(Line, "line:10%") {
		t.Error("Expected true, got false")
	}

	// Test case for Terminal pattern
	if !checkToken(Terminal, ".") {
		t.Error("Expected true, got false")
	}
}

func TestSearchTerminalToken(t *testing.T) {
	// Test case for token containing "."
	expected := []int{5}
	if result := searchTerminalToken("Hello."); !equalSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test case for token containing "?"
	expected = []int{5, 6}
	if result := searchTerminalToken("Hello?"); !equalSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test case for token not containing "." or "?"
	expected = nil
	if result := searchTerminalToken("Hello"); !equalSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCheckHeader(t *testing.T) {
	// Test case for matching "WEBVTT"
	if !checkHeader("WEBVTT") {
		t.Error("Expected true, got false")
	}

	// Test case for matching "Kind: captions"
	if !checkHeader("Kind: captions") {
		t.Error("Expected true, got false")
	}

	// Test case for not matching any header pattern
	if checkHeader("Hello") {
		t.Error("Expected false, got true")
	}
}

// equalSlices checks if two slices are equal.
func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
