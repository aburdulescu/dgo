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
	diagram, ast, err := Parse(strings.NewReader(string(content)))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(diagram)
	t.Log(ast)
}
