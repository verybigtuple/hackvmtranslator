package codewriter

import (
	"testing"

	"github.com/verybigtuple/hackvmtranslator/parser"
)

func TestWriterAdd(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdArithmeticBinary, Arg1: "add"}
	want := []string{
		"// add",
		"@SP",
		"AM=M-1",
		"D=M", // D = y

		"A=A-1",
		"M=D+M", // y + x
	}
	runTestLine(t, testLine, want)
}

func TestWriterSub(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdArithmeticBinary, Arg1: "sub"}
	want := []string{
		"// sub",
		"@SP",
		"AM=M-1",
		"D=M", // D = y

		"A=A-1",
		"M=M-D", // x-y
	}
	runTestLine(t, testLine, want)
}

func TestWriterNeg(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdArithmeticUnary, Arg1: "neg"}
	want := []string{
		"// neg",
		"@SP",
		"A=M-1",
		"M=-M",
	}
	runTestLine(t, testLine, want)
}

func TestWriterEq(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdArithmeticCond, Arg1: "eq"}
	want := []string{
		"// eq",
		"@SP",
		"AM=M-1",
		"D=M", // D = y

		"A=A-1",
		"D=M-D", // D=x-y
		"M=0",   // False by default

		"@test.EQ_END_0",
		"D;JNE", // If D=x-y != 0 then jump, esle set true (-1)

		"@SP",
		"A=M-1",
		"M=-1",
		"(test.EQ_END_0)",
	}
	runTestLine(t, testLine, want)
}

func TestWriterGt(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdArithmeticCond, Arg1: "gt"}
	want := []string{
		"// gt",
		"@SP",
		"AM=M-1",
		"D=M", // D = y

		"A=A-1",
		"D=M-D", // M=x; D=x-y. If (x > y), then x-y > 0
		"M=0",   // False by default

		"@test.GT_END_0",
		"D;JLE", // If D=x-y <= 0 then jump to end and leave M=False, else set true (-1)

		"@SP",
		"A=M-1",
		"M=-1",
		"(test.GT_END_0)",
	}
	runTestLine(t, testLine, want)
}

func TestWriterLt(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdArithmeticCond, Arg1: "lt"}
	want := []string{
		"// lt",
		"@SP",
		"AM=M-1",
		"D=M", // D = y

		"A=A-1",
		"D=M-D", // M=x; D=x-y. If (x > y), then x-y < 0
		"M=0",   // False by default

		"@test.LT_END_0",
		"D;JGE", // If D=x-y >= 0 then jump and leave M=False, esle set true (-1)

		"@SP",
		"A=M-1",
		"M=-1",
		"(test.LT_END_0)",
	}
	runTestLine(t, testLine, want)
}

func TestWriterAnd(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdArithmeticBinary, Arg1: "and"}
	want := []string{
		"// and",
		"@SP",
		"AM=M-1",
		"D=M", // D = y

		"A=A-1",
		"M=D&M", // x&y
	}
	runTestLine(t, testLine, want)
}

func TestWriterOr(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdArithmeticBinary, Arg1: "or"}
	want := []string{
		"// or",
		"@SP",
		"AM=M-1",
		"D=M", // D = y

		"A=A-1",
		"M=D|M", // x|y
	}
	runTestLine(t, testLine, want)
}

func TestWriterNot(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdArithmeticUnary, Arg1: "not"}
	want := []string{
		"// not",
		"@SP",
		"A=M-1",
		"M=!M",
	}
	runTestLine(t, testLine, want)
}
