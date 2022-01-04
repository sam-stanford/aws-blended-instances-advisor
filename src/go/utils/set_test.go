package utils

import "testing"

func TestStringSet(t *testing.T) {
	set := make(StringSet)

	if set.Contains("SHOULD NOT CONTAIN") {
		t.Fatalf("empty StringSet contains string")
	}

	set.Remove("test1")
	if set.Contains("test1") {
		t.Fatalf("StringSet contains string \"%s\" when it shouldn't", "test1")
	}

	set.Add("test2")
	if !set.Contains("test2") {
		t.Fatalf("StringSet does not contain string \"%s\" when it should", "test2")
	}

	set.Add("test2") // Add again
	if !set.Contains("test2") {
		t.Fatalf("StringSet does not contain string \"%s\" when it should", "test2")
	}

	set.Remove("test2")
	if set.Contains("test2") {
		t.Fatalf("StringSet contains string \"%s\" when it shouldn't", "test2")
	}

	set.Add("test2")
	if !set.Contains("test2") {
		t.Fatalf("StringSet does not contain string \"%s\" when it should", "test2")
	}

	set.Add("test3")
	if !set.Contains("test3") {
		t.Fatalf("StringSet does not contain string \"%s\" when it should", "test3")
	}

	set.Remove("test3")
	if set.Contains("test3") {
		t.Fatalf("StringSet contains string \"%s\" when it shouldn't", "test3")
	}
	if !set.Contains("test2") {
		t.Fatalf("StringSet does not contain string \"%s\" when it should", "test2")
	}
}
