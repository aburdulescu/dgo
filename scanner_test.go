package dgo

import (
	"os"
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	content, err := os.ReadFile("example.dgo")
	if err != nil {
		t.Fatal(err)
	}

	tokens := scan(strings.NewReader(string(content)))

	for _, tok := range tokens {
		t.Log(tok.text)
	}
}
