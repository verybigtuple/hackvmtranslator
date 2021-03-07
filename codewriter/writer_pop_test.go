package codewriter

import (
	"testing"

	"github.com/verybigtuple/hackvmtranslator/parser"
)

func TestWriterPopStatic(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdPop, Arg1: "static", Arg2: 5}
	want := []string{
		"// pop static 5",
		"@SP",
		"AM=M-1",
		"D=M",
		"@test.5",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPopLocalZero(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdPop, Arg1: "local", Arg2: 0}
	want := []string{
		"// pop local 0",
		"@SP",
		"AM=M-1",
		"D=M",
		"@LCL",
		"A=M",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPopLocalOne(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdPop, Arg1: "local", Arg2: 1}
	want := []string{
		"// pop local 1",
		"@SP",
		"AM=M-1",
		"D=M",
		"@LCL",
		"A=M+1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPopLocalTwo(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdPop, Arg1: "local", Arg2: 2}
	want := []string{
		"// pop local 2",
		"@SP",
		"AM=M-1",
		"D=M",
		"@LCL",
		"A=M+1",
		"A=A+1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPopLocalThree(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdPop, Arg1: "local", Arg2: 3}
	want := []string{
		"// pop local 3",
		"@SP",
		"AM=M-1",
		"D=M",
		"@LCL",
		"A=M+1",
		"A=A+1",
		"A=A+1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

// This case has 12 processor commands. For Arg2 > 7 all other case
func TestWriterPopLocalSeven(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdPop, Arg1: "local", Arg2: 7}
	want := []string{
		"// pop local 7",
		"@SP",
		"AM=M-1",
		"D=M",
		"@LCL",
		"A=M+1",
		"A=A+1",
		"A=A+1",
		"A=A+1",
		"A=A+1",
		"A=A+1",
		"A=A+1",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPopLocalMore(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdPop, Arg1: "local", Arg2: 8}
	want := []string{
		"// pop local 8",
		"@8",
		"D=A",
		"@LCL",
		"D=D+M",
		"@R13",
		"M=D",

		"@SP",
		"AM=M-1",
		"D=M",

		"@R13",
		"A=M",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPopPointerZero(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdPop, Arg1: "pointer", Arg2: 0}
	want := []string{
		"// pop pointer 0",
		"@SP",
		"AM=M-1",
		"D=M",
		"@THIS",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestWriterPopPointerOne(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdPop, Arg1: "pointer", Arg2: 1}
	want := []string{
		"// pop pointer 1",
		"@SP",
		"AM=M-1",
		"D=M",
		"@THAT",
		"M=D",
	}
	runTestLine(t, testLine, want)
}
