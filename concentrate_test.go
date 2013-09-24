package main

import (
	"testing"
)

func TestRemove(t *testing.T) {
	hf := HostsFile{Entries: []string{"www.reddit.com"}}
	hf.Remove("www.reddit.com")

	if hf.Entries[0] != "" {
		t.Fatal("Removed domain should be the empty string")
	}
}

func TestAddIdempotent(t *testing.T) {
	hf := HostsFile{Entries: []string{"www.reddit.com"}}
	hf.Add("www.reddit.com")

	if len(hf.Entries) != 1 {
		t.Fatal("Added a domain alread there")
	}
}

func TestAdd(t *testing.T) {
	hf := HostsFile{}
	hf.Add("www.reddit.com")

	if len(hf.Entries) != 1 {
		t.Fatal("Added a domain alread there")
	}
}
