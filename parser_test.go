package dgo

import (
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	content, err := os.ReadFile("example.dgo")
	if err != nil {
		t.Fatal(err)
	}
	err = parse(strings.NewReader(string(content)))
	if err != nil {
		t.Fatal(err)
	}
}
