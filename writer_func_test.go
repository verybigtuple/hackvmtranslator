package main

import (
	"testing"
)

func TestFuncFunction0(t *testing.T) {
	testLine := Command{CmdType: cmdFunction, Arg1: "Test.func", Arg2: 0}
	want := []string{
		"// function Test.func 0",
		"(Test.func)", // Push return address to stack
	}
	runTestLine(t, testLine, want)
}

func TestFuncFunction1(t *testing.T) {
	testLine := Command{CmdType: cmdFunction, Arg1: "Test.func", Arg2: 1}
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
	testLine := Command{CmdType: cmdFunction, Arg1: "Test.func", Arg2: 5}
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
	testLine := Command{CmdType: cmdCall, Arg1: "Test.func", Arg2: 2}
	want := []string{
		"// call Test.func 2",
		"@Test.func$return", // Push return address to stack
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

		"@THAT", // Push that to stack
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

		"(Test.func$return)",
	}
	runTestLine(t, testLine, want)
}

func TestFuncReturn(t *testing.T) {
	testLine := Command{CmdType: cmdReturn}
	want := []string{
		"// return",
		
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

		// Before changing *LCL, we should get *(Endframe-5)
		"@LCL",
		"A=M-1", // A = EndFrame - 4
		"A=A-1", // A = EndFrame - 5
		"D=M",   // D = *(EndFrame - 5)
		"@R14",
		"M=D", // RetrAddr = R14 = *(EndFrame - 5)

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
