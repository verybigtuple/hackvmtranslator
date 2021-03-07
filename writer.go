package main

import (
	"bufio"
	"fmt"
)

// CodeWriter is a struc that writes instructions to a user's writer
type CodeWriter struct {
	writer *bufio.Writer
	asm    *asmBuilder

	name      string
	stPrefix  string
	fnPrefix  string
	eqCount   int
	gtCount   int
	ltCount   int
	callCount int
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

var writers = map[CommandType]func(*CodeWriter, Command) (err error){
	cmdPush:             (*CodeWriter).writePush,
	cmdPop:              (*CodeWriter).writePop,
	cmdArithmeticBinary: (*CodeWriter).writeAritmBinary,
	cmdArithmeticUnary:  (*CodeWriter).writeArithmUnary,
	cmdArithmeticCond:   (*CodeWriter).writeArithmCond,
	cmdGoto:             (*CodeWriter).writeGotoCmd,
	cmdLabel:            (*CodeWriter).writeLabelCmd,
	cmdIfGoto:           (*CodeWriter).writeIfGotoCmd,
	cmdFunction:         (*CodeWriter).writeFunctionCmd,
	cmdCall:             (*CodeWriter).writeCallCmd,
	cmdReturn:           (*CodeWriter).writeReturnCmd,
}

// WriteCommand writes a command to a writer passed to NewCodeWriter
func (cw *CodeWriter) WriteCommand(cmd Command) error {
	if w, ok := writers[cmd.CmdType]; ok {
		return w(cw, cmd)
	}
	return fmt.Errorf("There is no writer for cmd")
}

func (cw *CodeWriter) WriteBootstrap() (err error) {
	// Init SP
	cw.asm.AsmCmds(256, "D=A", SP, "M=D")
	// Call Sys.init function
	cw.writeCallCmd(Command{cmdCall, "Sys.init", 0})
	// In order not to have 2 lablels in a row
	cw.asm.AsmCmds("D=0")

	_, err = cw.writer.WriteString(cw.asm.CodeAsm())
	return
}

func (cw *CodeWriter) writePush(cmd Command) error {
	cw.asm.AddComment(fmt.Sprintf("push %s %d", cmd.Arg1, cmd.Arg2))

	switch {
	case isConstantSegment(cmd.Arg1): // push constant 2
		cw.asm.AsmCmds(cmd.Arg2, "D=A")
	case isStaticSegment(cmd.Arg1): // push  static 2
		cw.asm.StaticAinstr(cw.stPrefix, cmd.Arg2)
		cw.asm.AsmCmds("D=M")
	case isTempSegment(cmd.Arg1): // push temp 2
		//cw.asm.TempToD(cmd.Arg2)
		cw.asm.TempAInstr(cmd.Arg2)
		cw.asm.AsmCmds("D=M")
	case isPointerSegment(cmd.Arg1):
		cw.asm.PointerAinstr(cmd.Arg2)
		cw.asm.AsmCmds("D=M")
	default: // push local 2
		if cmd.Arg2 <= 3 {
			cw.asm.SegmAddr(cmd.Arg1, cmd.Arg2)
		} else {
			cw.asm.SegmAddrCalcWithD(cmd.Arg1, cmd.Arg2, "A")
		}
		cw.asm.FromMemToD()
	}

	cw.asm.ToStack("D")
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writePop(cmd Command) error {
	cw.asm.AddComment(fmt.Sprintf("pop %s %d", cmd.Arg1, cmd.Arg2))

	switch {
	case isStaticSegment(cmd.Arg1):
		cw.asm.FromStack("D")
		cw.asm.StaticAinstr(cw.stPrefix, cmd.Arg2)
		cw.asm.AsmCmds("M=D")
	case isTempSegment(cmd.Arg1):
		cw.asm.FromStack("D")
		cw.asm.TempAInstr(cmd.Arg2)
		cw.asm.AsmCmds("M=D")
	case isPointerSegment(cmd.Arg1):
		cw.asm.FromStack("D")
		cw.asm.PointerAinstr(cmd.Arg2)
		cw.asm.AsmCmds("M=D")
	default:
		if cmd.Arg2 <= 7 {
			cw.asm.FromStack("D")
			cw.asm.SegmAddr(cmd.Arg1, cmd.Arg2)
		} else {
			cw.asm.SegmAddrCalcWithD(cmd.Arg1, cmd.Arg2, "D")
			cw.asm.ToR("D")
			cw.asm.FromStack("D")
			cw.asm.FromR("A")
		}
		cw.asm.FromDtoMem()
	}

	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeAritmBinary(cmd Command) error {
	cw.asm.AddComment(cmd.Arg1)
	cw.asm.FromStack("D")
	cw.asm.DecAddr()
	switch cmd.Arg1 {
	case "add":
		cw.asm.AsmCmds("M=D+M")
	case "sub":
		cw.asm.AsmCmds("M=M-D")
	case "and":
		cw.asm.AsmCmds("M=D&M")
	case "or":
		cw.asm.AsmCmds("M=D|M")
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
	cw.asm.FromStack("D")
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
	cw.asm.AsmCmds("0;JMP")
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
	cw.asm.FromStack("D")
	cw.asm.AtFuncLabel(cw.fnPrefix, cmd.Arg1)
	cw.asm.AsmCmds("D;JNE")
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeFunctionCmd(cmd Command) error {
	cw.fnPrefix = cmd.Arg1

	cw.asm.AddComment(fmt.Sprintf("function %s %d", cmd.Arg1, cmd.Arg2))
	cw.asm.SetLabel(cmd.Arg1)

	// If function has just one local var, then just push one zero to the stack
	if cmd.Arg2 == 1 {
		cw.asm.ToStack("0")
	}
	// If function has has 2 and more vars, then we can slightly oprimized initialization
	if cmd.Arg2 > 1 {
		// Init first local var to stack w/o moving SP pointer forward
		cw.asm.AsmCmds(SP, "A=M", "M=0")
		// Init the the rest of vars
		for i := 0; i < cmd.Arg2-1; i++ {
			cw.asm.AsmCmds("A=A+1", "M=0")
		}
		// Restore the right position in SP
		cw.asm.AsmCmds("D=A+1", "@SP", "M=D")
	}
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeCallCmd(cmd Command) error {
	cw.asm.AddComment(fmt.Sprintf("call %s %d", cmd.Arg1, cmd.Arg2))

	label := fmt.Sprintf("%s.CALL_RET_%d", cw.stPrefix, cw.callCount)
	cw.callCount++

	// Add redturnAddr to stack but do not move SP Pointer
	cw.asm.AtLabel(label)
	cw.asm.AsmCmds("D=A", SP, "A=M", "M=D")
	// Save all segments to the stack except for THAT (the last one)
	segm := [...]SegmInstr{LCL, ARG, THIS}
	for _, s := range segm {
		cw.asm.AsmCmds(s, "D=M", SP, "AM=M+1", "M=D")
	}
	// Save THAT to the stack and set SP Pointer to the normal value (empty stack register)
	cw.asm.AsmCmds(THAT, "D=M", SP, "M=M+1", "M=M+1", "A=M-1", "M=D")
	// Calc new ARG value - it is ARG = SP-5-<func args>
	offset := 5 + cmd.Arg2
	cw.asm.AsmCmds(offset, "D=A", SP, "D=M-D", ARG, "M=D")
	// Set new LCL value: LCL=SP
	cw.asm.AsmCmds(SP, "D=M", LCL, "M=D")
	// Jump to called function
	cw.asm.AtLabel(cmd.Arg1)
	cw.asm.AsmCmds("0;JMP")
	// Label of return address
	cw.asm.SetLabel(label)

	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeReturnCmd(cmd Command) error {
	cw.asm.AddComment("return")
	// Save return address. R14 = *(EndFrame - 5)
	cw.asm.AsmCmds(5, "D=A", LCL, "A=M-D", "D=M", "@R14", "M=D")
	// Move return value to arg. *ARG = Pop()
	cw.asm.FromStack("D")
	cw.asm.AsmCmds(ARG, "A=M", "M=D")
	// Recycle stack: SP = ARG + 1
	cw.asm.AsmCmds(ARG, "D=M+1", SP, "M=D")
	// Restore all func segments from the old stack
	segm := [...]SegmInstr{THAT, THIS, ARG, LCL}
	for _, s := range segm {
		cw.asm.AsmCmds(LCL, "AM=M-1", "D=M", s, "M=D")
	}
	// Jump to return address
	cw.asm.AsmCmds("@R14", "A=M", "0;JMP")

	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}
