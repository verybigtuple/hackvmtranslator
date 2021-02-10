package main

import (
	"bufio"
	"fmt"
)

// CodeWriter is a struc that writes instructions to a user's writer
type CodeWriter struct {
	writer *bufio.Writer
	asm    *asmBuilder

	stPrefix string
	eqCount  int
	gtCount  int
	ltCount  int
}

// NewCodeWriter retuns a pointer to a new CodeWriter
func NewCodeWriter(stPr string, w *bufio.Writer) *CodeWriter {
	return &CodeWriter{stPrefix: stPr, writer: w, asm: newAsmBuilder()}
}

// WriteCommand writes a command to a writer passed to NewCodeWriter
func (cw *CodeWriter) WriteCommand(cmd Command) (err error) {
	switch cmd.CmdType {
	case cPush:
		err = cw.writePush(cmd)
	case cPop:
		err = cw.writePop(cmd)
	case cArithmeticBinary:
		err = cw.writeAritmBinary(cmd)
	case cArithmeticUnary:
		err = cw.writeArithmUnary(cmd)
	case cArithmeticCond:
		err = cw.writeArithmCond(cmd)
	}
	return
}

func (cw *CodeWriter) writePush(cmd Command) error {
	cw.asm.AddComment(fmt.Sprintf("push %s %d", cmd.Arg1, cmd.Arg2))

	switch {
	case isConstantSegment(cmd.Arg1): // push constant 2
		cw.asm.ConstToD(cmd.Arg2)
	case isStaticSegment(cmd.Arg1): // push  static 2
		cw.asm.StaticToD(cw.stPrefix, cmd.Arg2)
	case isTempSegment(cmd.Arg1): // push temp 2
		cw.asm.TempToD(cmd.Arg2)
	default: // push local 2
		if cmd.Arg2 <= 3 {
			cw.asm.SegmAddr(cmd.Arg1, cmd.Arg2)
		} else {
			cw.asm.SegmAddrCalcWithD(cmd.Arg1, cmd.Arg2, "A")
		}
		cw.asm.FromMemToD()
	}

	cw.asm.FromDtoStack()
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writePop(cmd Command) error {
	cw.asm.AddComment(fmt.Sprintf("pop %s %d", cmd.Arg1, cmd.Arg2))

	switch {
	case isStaticSegment(cmd.Arg1):
		cw.asm.FromStackToD()
		cw.asm.StaticFromD(cw.stPrefix, cmd.Arg2)
	case isTempSegment(cmd.Arg1):
		cw.asm.FromStackToD()
		cw.asm.TempFromD(cmd.Arg2)
	default:
		if cmd.Arg2 <= 7 {
			cw.asm.FromStackToD()
			cw.asm.SegmAddr(cmd.Arg1, cmd.Arg2)
		} else {
			cw.asm.SegmAddrCalcWithD(cmd.Arg1, cmd.Arg2, "D")
			cw.asm.ToR("D")
			cw.asm.FromStackToD()
			cw.asm.FromR("A")
		}
		cw.asm.FromDtoMem()
	}

	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeAritmBinary(cmd Command) error {
	cw.asm.AddComment(cmd.Arg1)
	cw.asm.FromStackToD()
	cw.asm.DecAddr()
	switch cmd.Arg1 {
	case "add":
		cw.asm.ArbitraryCmd("M=D+M")
	case "sub":
		cw.asm.ArbitraryCmd("M=M-D")
	case "and":
		cw.asm.ArbitraryCmd("M=D&M")
	case "or":
		cw.asm.ArbitraryCmd("M=D|M")
	}
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeArithmUnary(cmd Command) error {
	cw.asm.AddComment(cmd.Arg1)
	switch cmd.Arg1 {
	case "neg":
		cw.asm.SetTopStack("-M")
	case "not":
		cw.asm.SetTopStack("!M")
	}
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeArithmCond(cmd Command) error {
	cw.asm.AddComment(cmd.Arg1)
	cw.asm.FromStackToD()
	cw.asm.CondFalseDefault()
	switch cmd.Arg1 {
	case "eq":
		cw.asm.CondJump(cmd.Arg1, cw.eqCount, "JNE")
		cw.eqCount++
	case "gt":
		cw.asm.CondJump(cmd.Arg1, cw.gtCount, "JLE")
		cw.gtCount++
	case "lt":
		cw.asm.CondJump(cmd.Arg1, cw.ltCount, "JGE")
		cw.ltCount++
	}

	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}
