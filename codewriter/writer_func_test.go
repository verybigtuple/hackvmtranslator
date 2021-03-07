package codewriter

import (
	"testing"

	"github.com/verybigtuple/hackvmtranslator/parser"
)

func TestFuncFunction0(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdFunction, Arg1: "Test.func", Arg2: 0}
	want := []string{
		"// function Test.func 0",
		"(Test.func)", // Push return address to stack
	}
	runTestLine(t, testLine, want)
}

func TestFuncFunction1(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdFunction, Arg1: "Test.func", Arg2: 1}
	want := []string{
		"// function Test.func 1",
		"(Test.func)", // Push return address to stack
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=0",
	}
	runTestLine(t, testLine, want)
}

func TestFuncFunctionMore(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdFunction, Arg1: "Test.func", Arg2: 5}
	want := []string{
		"// function Test.func 5",
		"(Test.func)", // Push return address to stack
		"@SP",
		"A=M",
		"M=0", // 1
		"A=A+1",
		"M=0", // 2
		"A=A+1",
		"M=0", // 3
		"A=A+1",
		"M=0", // 4
		"A=A+1",
		"M=0",   // 5
		"D=A+1", // D = SP
		"@SP",
		"M=D",
	}
	runTestLine(t, testLine, want)
}

func TestFuncCall(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdCall, Arg1: "Test.func", Arg2: 2}
	want := []string{
		"// call Test.func 2",
		"@test.CALL_RET_0", // Push return address to stack
		"D=A",
		"@SP",
		"A=M",
		"M=D", // Push without incrementing SP, as SP is ponting to an empty register

		"@LCL", // Push local to stack
		"D=M",
		"@SP",
		"AM=M+1",
		"M=D",

		"@ARG", // Push arg to stack
		"D=M",
		"@SP",
		"AM=M+1",
		"M=D",

		"@THIS", // Push this to stack
		"D=M",
		"@SP",
		"AM=M+1",
		"M=D",

		"@THAT", // Push that to stack and restore SP to the normal position
		"D=M",
		"@SP",
		"M=M+1",
		"M=M+1",
		"A=M-1",
		"M=D",

		"@7", // 5+2
		"D=A",
		"@SP",
		"D=M-D", // Arg = SP - 5 - nargs
		"@ARG",
		"M=D",

		"@SP", // LCL = SP
		"D=M",
		"@LCL",
		"M=D",

		"@Test.func",
		"0;JMP",

		"(test.CALL_RET_0)",
	}
	runTestLine(t, testLine, want)
}

func TestFuncReturn(t *testing.T) {
	testLine := parser.Command{CmdType: parser.CmdReturn}
	want := []string{
		"// return",

		// Save return address to R14
		// We cannot do it later, Since if function is called with zero arguments,
		// ARG will point to the same RAM address where ReturnAddress contains.
		// In this case *ARG = Pop() will erase return address with return value
		"@5",
		"D=A",
		"@LCL",
		"A=M-D",
		"D=M",
		"@R14",
		"M=D",

		// *ARG = Pop()
		"@SP",
		"AM=M-1",
		"D=M", // D = Pop()
		"@ARG",
		"A=M",
		"M=D", //  *ARG = Pop()

		// SP = ARG + 1
		"@ARG",
		"D=M+1",
		"@SP",
		"M=D", // Recycle stack

		// LCL
		"@LCL",   // EndFrame = LCL
		"AM=M-1", // A = Endframe-1, LCL = Endframe-1
		"D=M",    // D = *(Endframe-1)
		"@THAT",
		"M=D",

		"@LCL",
		"AM=M-1", // A = Endframe-2, LCL = Endframe-2
		"D=M",    // D = *(Endframe-2)
		"@THIS",
		"M=D",

		"@LCL",
		"AM=M-1", // A = Endframe-3, LCL = Endframe-3
		"D=M",    // D = *(Endframe-3)
		"@ARG",
		"M=D",

		"@LCL",
		"AM=M-1", // A = Endframe-4, LCL = Endframe-4
		"D=M",
		"@LCL",
		"M=D",

		// Goto RetrAddress
		"@R14",
		"A=M",
		"0;JMP",
	}
	runTestLine(t, testLine, want)
}
