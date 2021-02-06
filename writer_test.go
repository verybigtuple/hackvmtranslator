package main

import (
	"bufio"
	"strings"
	"testing"
)

func runTestLine(t *testing.T, tc Command, want []string) {
	sb := strings.Builder{}
	writer := bufio.NewWriter(&sb)

	codeWriter := NewCodeWriter("test", writer)
	err := codeWriter.WriteCommand(tc)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	writer.Flush()

	splitFunc := func(r rune) bool {
		return r == '\n'
	}

	actual := strings.FieldsFunc(sb.String(), splitFunc)

	if len(actual) != len(want) {
		t.Errorf("Actual len: %v; want len: %v", len(actual), len(want))
		return
	}

	for i, actualLine := range actual {
		if actualLine != want[i] {
			t.Errorf("Line %v, Actual %v; want %v", i, actualLine, want[i])
			return
		}
	}
}

func TestWriterPushConst(t *testing.T) {
	testLine := Command{CmdType: cPush, Arg1: "constant", Arg2: 100}
	want := []string{
		"@100",
		"D=A",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPushStatic(t *testing.T) {
	testLine := Command{CmdType: cPush, Arg1: "static", Arg2: 5}
	want := []string{
		"@test.5",
		"D=M",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}
