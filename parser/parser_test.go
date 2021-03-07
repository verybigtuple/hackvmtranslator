package parser

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
			want: Command{CmdPush, "constant", 17},
		},
		{
			line: "pop local 1",
			want: Command{CmdPop, "local", 1},
		},
		{
			line: "add",
			want: Command{CmdArithmeticBinary, "add", 0},
		},
		{
			line: "push local 100 // Comment for the command",
			want: Command{CmdPush, "local", 100},
		},
		{
			line: "eq",
			want: Command{CmdArithmeticCond, "eq", 0},
		},
		{
			line: "goto testLabel",
			want: Command{CmdGoto, "testLabel", 0},
		},
		{
			line: "label testLabel",
			want: Command{CmdLabel, "testLabel", 0},
		},
		{
			line: "if-goto testLabel",
			want: Command{CmdIfGoto, "testLabel", 0},
		},
		{
			line: "function Main.test 2",
			want: Command{CmdFunction, "Main.test", 2},
		},
		{
			line: "call Main.test 2",
			want: Command{CmdCall, "Main.test", 2},
		},
		{
			line: "return",
			want: Command{CmdReturn, "", 0},
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

func TestParserErrors(t *testing.T) {
	testCase := []struct {
		desc string
		line string
	}{
		{"No necessary argument", "push"},
		{"No second arg", "pop local"},
		{"Exceeded args", "pop local 2 3"},
		{"Wrong push literal", "pushd local 2"},
		{"Wrong push segment", "push lcl 2"},
		{"Wrong arithmetic arg", "add local"},
		{"Too few goto arg", "label"},
		{"Too many goto arg", "goto label1 label2"},
		{"No offset in function", "function Main.Test"},
		{"Illegal offset in function", "function Main.Test -1"},
		{"Too many function args", "function Main.Test 0 1"},
		{"No offset in call", "call Main.Test"},
		{"Illegal offset in call", "call Main.Test -1"},
		{"Too many call args", "call Main.Test 0 1"},
		{"Args in return", "return 0"},
	}

	for _, tc := range testCase {
		t.Run(tc.desc, func(t *testing.T) {
			parser := newParserString(tc.line)
			cmd, err := parser.ParseNext()
			if err == nil {
				t.Errorf("Error is not arisen. Cmd %+v", *cmd)
				return
			}
			if errors.Is(err, io.EOF) {
				t.Errorf("Unexpected EOF error")
			}
		})
	}
}
