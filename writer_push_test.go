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
		"// push constant 100",
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
		"// push static 5",
		"@test.5",
		"D=M",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPushTempZero(t *testing.T) {
	testLine := Command{CmdType: cPush, Arg1: "temp", Arg2: 0}
	want := []string{
		"// push temp 0",
		"@5",
		"D=M",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPushTempOne(t *testing.T) {
	testLine := Command{CmdType: cPush, Arg1: "temp", Arg2: 1}
	want := []string{
		"// push temp 1",
		"@6", // 5+1
		"D=M",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPushLocalZero(t *testing.T) {
	testLine := Command{CmdType: cPush, Arg1: "local", Arg2: 0}
	want := []string{
		"// push local 0",
		"@LCL",
		"A=M",
		"D=M",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPushLocalOne(t *testing.T) {
	testLine := Command{CmdType: cPush, Arg1: "local", Arg2: 1}
	want := []string{
		"// push local 1",
		"@LCL",
		"A=M+1",
		"D=M",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPushLocalTwo(t *testing.T) {
	testLine := Command{CmdType: cPush, Arg1: "local", Arg2: 3}
	want := []string{
		"// push local 3",
		"@LCL",
		"A=M+1",
		"A=A+1",
		"A=A+1", // Optimization for push local 3 - 9 cmd
		"D=M",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPushLocalMore(t *testing.T) {
	testLine := Command{CmdType: cPush, Arg1: "local", Arg2: 5}
	want := []string{
		"// push local 5",
		"@5",
		"D=A",
		"@LCL",
		"A=D+M", // Calc lcl + 5
		"D=M",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}
