package dgo

import (
	"os"
	"strings"
	"testing"
)

func TestDiagram(t *testing.T) {
	content, err := os.ReadFile("example.dgo")
	if err != nil {
		t.Fatal(err)
	}
	ast, err := Parse(strings.NewReader(string(content)))
	if err != nil {
		t.Fatal(err)
	}

	d := NewDiagram(ast)
	t.Log(d)
}
