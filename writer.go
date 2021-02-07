package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
)

const tempBaseAddr = 5 // Base address for temp vars

// Assign segments to assembler A-Instructions
var segmAInstr = map[string]string{
	"local":    "@LCL",
	"argument": "@ARG",
	"this":     "@THIS",
	"that":     "@THAT",
}

// CodeWriter is a struc that writes instructions to a user's writer
type CodeWriter struct {
	stPrefix string
	writer   *bufio.Writer
	asm      *asmHelper
}

// NewCodeWriter retuns a pointer to a new CodeWriter
func NewCodeWriter(stPr string, w *bufio.Writer) *CodeWriter {
	return &CodeWriter{stPrefix: stPr, writer: w, asm: newAsmHelper()}
}

// WriteCommand writes a command to a writer passed to NewCodeWriter
func (cw *CodeWriter) WriteCommand(cmd Command) (err error) {
	switch cmd.CmdType {
	case cPush:
		err = cw.writePush(cmd)
	}
	return
}

func (cw *CodeWriter) writePush(cmd Command) (err error) {
	comment := fmt.Sprintf("%s push %s %d\n", commentPrefix, cmd.Arg1, cmd.Arg2)
	var code string

	switch {
	case isConstantSegment(cmd.Arg1): // push constant 2
		code = cw.asm.ConstToDReg(cmd.Arg2)
	case isStaticSegment(cmd.Arg1): // push  static 2
		code = cw.asm.StaticToDReg(cw.stPrefix, cmd.Arg2)
	case isTempSegment(cmd.Arg1): // push temp 2
		code = cw.asm.TempToReg(cmd.Arg2)
	default: // push local 2
		code = cw.asm.SegmentToDReg(cmd.Arg1, cmd.Arg2)
	}

	_, err = cw.writer.WriteString(comment + code + pushDReg)
	return
}

// Push value from the D-Register to the stack
var pushDReg = `
@SP
M=M+1
A=M-1
M=D
`

// Pop value from the stack and save it to the D-Register
var popDReg = `
@SP
AM=M-1
D=M
`

type asmHelper struct {
	builder *bytes.Buffer // byte buffer does not reallocates while reset (strings.Builder does)
}

func newAsmHelper() *asmHelper {
	b := bytes.Buffer{}
	return &asmHelper{&b}
}

func (ah *asmHelper) ConstToDReg(c int) string {
	ah.builder.Reset()
	ah.builder.WriteString("@" + strconv.Itoa(c) + "\n") // like @101
	ah.builder.WriteString("D=A\n")
	return ah.builder.String()
}

func (ah *asmHelper) SegmentToDReg(segm string, offset int) string {
	ah.builder.Reset()

	if offset > 3 {
		ah.builder.WriteString("@" + strconv.Itoa(offset) + "\n") // like @101
		ah.builder.WriteString("D=A\n")
	}

	ah.builder.WriteString(segmAInstr[segm] + "\n") // like @ARG
	switch {
	case offset == 0:
		ah.builder.WriteString("A=M\n")
	case offset == 1:
		ah.builder.WriteString("A=M+1\n")
	case offset >= 2 && offset <= 3:
		ah.builder.WriteString("A=M+1\n")
		for i := 0; i < offset-1; i++ {
			ah.builder.WriteString("A=A+1\n")
		}
	default:
		ah.builder.WriteString("A=D+M\n")
	}
	ah.builder.WriteString("D=M\n")
	return ah.builder.String()
}

func (ah *asmHelper) StaticToDReg(prefix string, id int) string {
	ah.builder.Reset()
	ah.builder.WriteString(fmt.Sprintf("@%s.%d\n", prefix, id)) //like @file.1
	ah.builder.WriteString("D=M")
	return ah.builder.String()
}

func (ah *asmHelper) TempToReg(offset int) string {
	ah.builder.Reset()
	addr := tempBaseAddr + offset
	ah.builder.WriteString("@" + strconv.Itoa(addr) + "\n")
	ah.builder.WriteString("D=M\n")
	return ah.builder.String()
}
