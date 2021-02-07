package main

import (
	"bufio"
	"fmt"
)

// CodeWriter is a struc that writes instructions to a user's writer
type CodeWriter struct {
	stPrefix string
	writer   *bufio.Writer
	asm      *asmBuilder
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
	}
	return
}

func (cw *CodeWriter) writePush(cmd Command) error {
	cw.asm.AddComment(fmt.Sprintf("push %s %d\n", cmd.Arg1, cmd.Arg2))

	switch {
	case isConstantSegment(cmd.Arg1): // push constant 2
		cw.asm.AddConstToDReg(cmd.Arg2)
	case isStaticSegment(cmd.Arg1): // push  static 2
		cw.asm.AddStaticToDReg(cw.stPrefix, cmd.Arg2)
	case isTempSegment(cmd.Arg1): // push temp 2
		cw.asm.AddTempToDReg(cmd.Arg2)
	default: // push local 2
		if cmd.Arg2 <= 3 {
			cw.asm.AddCalcSegmentAddr(cmd.Arg1, cmd.Arg2)
		} else {
			cw.asm.AddCalcSegmentAddrWithD(cmd.Arg1, cmd.Arg2, "A")
		}
		cw.asm.AddDeqM()
	}

	cw.asm.AddPushFromDReg()
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writePop(cmd Command) error {
	cw.asm.AddComment(fmt.Sprintf("pop %s %d\n", cmd.Arg1, cmd.Arg2))

	switch {
	case isStaticSegment(cmd.Arg1):
		cw.asm.AddPopToDReg()
		cw.asm.AddStaticFromDReg(cw.stPrefix, cmd.Arg2)
	case isTempSegment(cmd.Arg1):
		cw.asm.AddPopToDReg()
		cw.asm.AddTempFromDReg(cmd.Arg2)
	default:
		if cmd.Arg2 <= 7 {
			cw.asm.AddPopToDReg()
			cw.asm.AddCalcSegmentAddr(cmd.Arg1, cmd.Arg2)
		} else {
			cw.asm.AddCalcSegmentAddrWithD(cmd.Arg1, cmd.Arg2, "D")
			cw.asm.AddToRreg("D")
			cw.asm.AddPopToDReg()
			cw.asm.AddFromRreg("A")
		}
		cw.asm.AddMeqD()
	}

	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}
