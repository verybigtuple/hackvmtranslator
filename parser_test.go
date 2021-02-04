package main

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"
)

func newParserString(s string) *Parser {
	reader := bufio.NewReader(strings.NewReader(s))
	return NewParser(reader)
}

func TestParserRegular(t *testing.T) {
	testCases := []struct {
		line string
		want Command
	}{
		{
			line: "push constant 17",
			want: Command{cPush, "constant", 17},
		},
		{
			line: "pop local 1",
			want: Command{cPop, "local", 1},
		},
		{
			line: "add",
			want: Command{cArithmetic, "add", 0},
		},
		{
			line: "push local 100 // Comment for the command",
			want: Command{cPush, "local", 100},
		},
		{
			line: "eq",
			want: Command{cArithmetic, "eq", 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.line, func(t *testing.T) {
			parser := newParserString(tc.line)
			cmd, err := parser.ParseNext()
			if err != nil {
				t.Errorf("An error was returned: %v", err)
				return
			}
			if *cmd != tc.want {
				t.Errorf("actual: %+v; want: %+v", *cmd, tc.want)
			}
		})
	}
}

func TestParserMultiline(t *testing.T) {
	testCase := `
	// Some comment 


	push constant 100
	// Some other comment 

	pop local 0
	// end comment
	`
	parser := newParserString(testCase)
	result := make([]Command, 0, 2)

	for {
		cmd, err := parser.ParseNext()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Errorf("Some unexpected error arose: %v", err)
			return
		}
		result = append(result, *cmd)
	}
	if len(result) != 2 {
		t.Errorf("Got %v commands: %+v; Expected only 2 of them", len(result), result)
	}
}
