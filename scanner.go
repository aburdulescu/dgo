package dgo

import (
	"fmt"
	"io"
	"text/scanner"
)

const scannerMode = 0 |
	scanner.ScanIdents |
	scanner.ScanStrings |
	scanner.ScanRawStrings |
	scanner.ScanComments |
	scanner.SkipComments

type token struct {
	text   string
	line   uint32
	column uint32
}

func (t token) errorf(err error) error {
	if t.empty() {
		return err
	}
	return fmt.Errorf(
		"line=%d, column=%d, token='%s': %w",
		t.line, t.column, t.text, err,
	)
}

func (t token) empty() bool {
	return t.line == 0 && t.column == 0 && t.text == ""
}

func scan(r io.Reader) []token {
	var s scanner.Scanner

	s.Init(r)

	s.Mode = scannerMode

	var tokens []token
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		tokens = append(tokens, token{
			line:   uint32(s.Position.Line),
			column: uint32(s.Position.Column),
			text:   s.TokenText(),
		})
	}

	// add empty token to signify the end of file
	tokens = append(tokens, token{})

	return tokens
}
