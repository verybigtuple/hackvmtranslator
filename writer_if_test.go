package main

import "testing"

func TestWriterGoto(t *testing.T) {
	testLine := Command{CmdType: cmdGoto, Arg1: "label1"}
	want := []string{
		"// goto label1",
		"@func$label1",
		"0;JMP",
	}
	runTestLine(t, testLine, want)
}

func TestWriterLablel(t *testing.T) {
	testLine := Command{CmdType: cmdLabel, Arg1: "label1"}
	want := []string{
		"// label label1",
		"(func$label1)",
	}
	runTestLine(t, testLine, want)
}

func TestWriterIfGoto(t *testing.T) {
	testLine := Command{CmdType: cmdIfGoto, Arg1: "label1"}
	want := []string{
		"// if-goto label1",
		"@SP",
		"AM=M-1",
		"D=M",
		"@func$label1",
		"D;JNE",
	}
	runTestLine(t, testLine, want)
}
