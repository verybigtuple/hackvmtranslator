package main

import (
	"bufio"
	"fmt"
)

// CodeWriter is a struc that writes instructions to a user's writer
type CodeWriter struct {
	writer *bufio.Writer
	asm    *asmBuilder

	name     string
	stPrefix string
	fnPrefix string
	eqCount  int
	gtCount  int
	ltCount  int
}

// NewCodeWriter retuns a pointer to a new CodeWriter
func NewCodeWriter(w *bufio.Writer, name, stPrefix, fnPrefix string) *CodeWriter {
	if fnPrefix == "" {
		fnPrefix = "default"
	}
	cw := CodeWriter{
		writer:   w,
		asm:      newAsmBuilder(),
		name:     name,
		stPrefix: stPrefix,
		fnPrefix: fnPrefix,
	}

	if name != "" {
		cw.asm.AddComment(name)
	}
	return &cw
}

// NewCodeWriterBootstrap creates Codewriter for Bootstrap
func NewCodeWriterBootstrap(w *bufio.Writer) *CodeWriter {
	return NewCodeWriter(w, "Bootstrap", "", "")
}

// WriteCommand writes a command to a writer passed to NewCodeWriter
func (cw *CodeWriter) WriteCommand(cmd Command) (err error) {
	switch cmd.CmdType {
	case cmdPush:
		err = cw.writePush(cmd)
	case cmdPop:
		err = cw.writePop(cmd)
	case cmdArithmeticBinary:
		err = cw.writeAritmBinary(cmd)
	case cmdArithmeticUnary:
		err = cw.writeArithmUnary(cmd)
	case cmdArithmeticCond:
		err = cw.writeArithmCond(cmd)
	case cmdGoto:
		err = cw.writeGotoCmd(cmd)
	case cmdLabel:
		err = cw.writeLabelCmd(cmd)
	case cmdIfGoto:
		err = cw.writeIfGotoCmd(cmd)
	}
	return
}

func (cw *CodeWriter) WriteBootstrap() (err error) {
	// Init SP
	cw.asm.ArbitraryCmd("@256")
	cw.asm.ArbitraryCmd("D=A")
	cw.asm.ArbitraryCmd("@SP")
	cw.asm.ArbitraryCmd("M=D")
	// Go to Sys.init function
	cw.asm.AtLabel("Sys.init")
	cw.asm.ArbitraryCmd("0;JMP")

	_, err = cw.writer.WriteString(cw.asm.CodeAsm())
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
	case isPointerSegment(cmd.Arg1):
		cw.asm.PointerToD(cmd.Arg2)
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
	case isPointerSegment(cmd.Arg1):
		cw.asm.FromStackToD()
		cw.asm.PointerFromD(cmd.Arg2)
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
		cw.asm.CondJump(cw.stPrefix, cmd.Arg1, "JNE", cw.eqCount)
		cw.eqCount++
	case "gt":
		cw.asm.CondJump(cw.stPrefix, cmd.Arg1, "JLE", cw.gtCount)
		cw.gtCount++
	case "lt":
		cw.asm.CondJump(cw.stPrefix, cmd.Arg1, "JGE", cw.ltCount)
		cw.ltCount++
	}

	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeGotoCmd(cmd Command) error {
	cw.asm.AddComment("goto " + cmd.Arg1)
	cw.asm.AtFuncLabel(cw.fnPrefix, cmd.Arg1)
	cw.asm.ArbitraryCmd("0;JMP")
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeLabelCmd(cmd Command) error {
	cw.asm.AddComment("label " + cmd.Arg1)
	cw.asm.SetFuncLabel(cw.fnPrefix, cmd.Arg1)
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeIfGotoCmd(cmd Command) error {
	cw.asm.AddComment("if-goto " + cmd.Arg1)
	cw.asm.FromStackToD()
	cw.asm.AtFuncLabel(cw.fnPrefix, cmd.Arg1)
	cw.asm.ArbitraryCmd("D;JNE")
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}
