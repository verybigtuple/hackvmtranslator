package codewriter

import (
	"bufio"
	"fmt"

	"github.com/verybigtuple/hackvmtranslator/parser"
)

// CodeWriter is a struc that writes instructions to a user's writer
type CodeWriter struct {
	writer *bufio.Writer
	asm    *asmBuilder

	name        string
	stPrefix    string
	fnPrefix    string
	arCondCount int
	callCount   int
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

var writers = map[parser.CommandType]func(*CodeWriter, parser.Command) (err error){
	parser.CmdPush:             (*CodeWriter).writePush,
	parser.CmdPop:              (*CodeWriter).writePop,
	parser.CmdArithmeticBinary: (*CodeWriter).writeAritmBinary,
	parser.CmdArithmeticUnary:  (*CodeWriter).writeArithmUnary,
	parser.CmdArithmeticCond:   (*CodeWriter).writeArithmCond,
	parser.CmdGoto:             (*CodeWriter).writeGotoCmd,
	parser.CmdLabel:            (*CodeWriter).writeLabelCmd,
	parser.CmdIfGoto:           (*CodeWriter).writeIfGotoCmd,
	parser.CmdFunction:         (*CodeWriter).writeFunctionCmd,
	parser.CmdCall:             (*CodeWriter).writeCallCmd,
	parser.CmdReturn:           (*CodeWriter).writeReturnCmd,
}

// WriteCommand writes a command to a writer passed to NewCodeWriter
func (cw *CodeWriter) WriteCommand(cmd parser.Command) error {
	if w, ok := writers[cmd.CmdType]; ok {
		return w(cw, cmd)
	}
	return fmt.Errorf("There is no writer for cmd")
}

func (cw *CodeWriter) WriteBootstrap() (err error) {
	// Init SP
	cw.asm.AsmCmds(256, "D=A", sp, "M=D")
	// Call Sys.init function
	cw.writeCallCmd(parser.Command{parser.CmdCall, "Sys.init", 0})
	// In order not to have 2 lablels in a row
	cw.asm.AsmCmds("D=0")

	_, err = cw.writer.WriteString(cw.asm.CodeAsm())
	return
}

func (cw *CodeWriter) writePush(cmd parser.Command) error {
	cw.asm.AddComment(fmt.Sprintf("push %s %d", cmd.Arg1, cmd.Arg2))

	switch {
	case parser.IsConstantSegment(cmd.Arg1): // push constant 2
		cw.asm.AsmCmds(cmd.Arg2, "D=A")
	case parser.IsStaticSegment(cmd.Arg1): // push  static 2
		cw.asm.StaticAinstr(cw.stPrefix, cmd.Arg2)
		cw.asm.AsmCmds("D=M")
	case parser.IsTempSegment(cmd.Arg1): // push temp 2
		//cw.asm.TempToD(cmd.Arg2)
		cw.asm.TempAInstr(cmd.Arg2)
		cw.asm.AsmCmds("D=M")
	case parser.IsPointerSegment(cmd.Arg1):
		cw.asm.PointerAinstr(cmd.Arg2)
		cw.asm.AsmCmds("D=M")
	default: // push local, push argument or this/that
		// If offset is <=3 some optimisation is possible
		if cmd.Arg2 <= 3 {
			cw.asm.SegmentAinstr(cmd.Arg1)
			if cmd.Arg2 == 0 {
				cw.asm.AsmCmds("A=M")
			} else {
				cw.asm.AsmCmds("A=M+1")
			}
			for i := 0; i < cmd.Arg2-1; i++ {
				cw.asm.AsmCmds("A=A+1")
			}
		} else {
			cw.asm.AsmCmds(cmd.Arg2, "D=A")
			cw.asm.SegmentAinstr(cmd.Arg1)
			cw.asm.AsmCmds("A=D+M")
		}
		cw.asm.AsmCmds("D=M")
	}

	cw.asm.ToStack("D")
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writePop(cmd parser.Command) error {
	cw.asm.AddComment(fmt.Sprintf("pop %s %d", cmd.Arg1, cmd.Arg2))

	switch {
	case parser.IsStaticSegment(cmd.Arg1):
		cw.asm.FromStack("D")
		cw.asm.StaticAinstr(cw.stPrefix, cmd.Arg2)
		cw.asm.AsmCmds("M=D")
	case parser.IsTempSegment(cmd.Arg1):
		cw.asm.FromStack("D")
		cw.asm.TempAInstr(cmd.Arg2)
		cw.asm.AsmCmds("M=D")
	case parser.IsPointerSegment(cmd.Arg1):
		cw.asm.FromStack("D")
		cw.asm.PointerAinstr(cmd.Arg2)
		cw.asm.AsmCmds("M=D")
	default:
		if cmd.Arg2 <= 7 {
			cw.asm.FromStack("D")
			cw.asm.SegmentAinstr(cmd.Arg1)
			if cmd.Arg2 == 0 {
				cw.asm.AsmCmds("A=M")
			} else {
				cw.asm.AsmCmds("A=M+1")
			}
			for i := 0; i < cmd.Arg2-1; i++ {
				cw.asm.AsmCmds("A=A+1")
			}
		} else {
			cw.asm.AsmCmds(cmd.Arg2, "D=A")
			cw.asm.SegmentAinstr(cmd.Arg1)
			cw.asm.AsmCmds("D=D+M", r13, "M=D")
			cw.asm.FromStack("D")
			cw.asm.AsmCmds(r13, "A=M")
		}
		cw.asm.AsmCmds("M=D")
	}
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeAritmBinary(cmd parser.Command) error {
	cw.asm.AddComment(cmd.Arg1)
	cw.asm.FromStack("D")
	cw.asm.AsmCmds("A=A-1")
	switch cmd.Arg1 {
	case parser.AddKey:
		cw.asm.AsmCmds("M=D+M")
	case parser.SubKey:
		cw.asm.AsmCmds("M=M-D")
	case parser.AndKey:
		cw.asm.AsmCmds("M=D&M")
	case parser.OrKey:
		cw.asm.AsmCmds("M=D|M")
	}
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeArithmUnary(cmd parser.Command) error {
	cw.asm.AddComment(cmd.Arg1)
	// Get address for result (top of the stack)
	cw.asm.AsmCmds(sp, "A=M-1")
	// make calculation
	switch cmd.Arg1 {
	case parser.NegKey:
		cw.asm.AsmCmds("M=-M")
	case parser.NotKey:
		cw.asm.AsmCmds("M=!M")
	}
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeArithmCond(cmd parser.Command) error {
	cw.asm.AddComment(cmd.Arg1)
	// Get boolean from stack to D-register
	cw.asm.FromStack("D")
	// By default set to false
	cw.asm.AsmCmds("A=A-1", "D=M-D", "M=0")
	// Set label to jump if condition is true
	cw.asm.AtArithmCondLabel(cw.stPrefix, cmd.Arg1, cw.arCondCount)

	switch cmd.Arg1 {
	case parser.EqKey:
		cw.asm.AsmCmds("D;JNE") // if D=M-D != 0 than jump to the end and leave M=false
	case parser.GtKey:
		cw.asm.AsmCmds("D;JLE") // if D=M-D <=0 then jump to the end and leave M=false
	case parser.LtKey:
		cw.asm.AsmCmds("D;JGE") // if D=M-D >= 0 then jump to the end and leave M=false
	}
	// Set true
	cw.asm.AsmCmds(sp, "A=M-1", "M=-1")
	cw.asm.SetArithmCondLabel(cw.stPrefix, cmd.Arg1, cw.arCondCount)
	cw.arCondCount++
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeGotoCmd(cmd parser.Command) error {
	cw.asm.AddComment("goto " + cmd.Arg1)
	cw.asm.AtFuncLabel(cw.fnPrefix, cmd.Arg1)
	cw.asm.AsmCmds("0;JMP")
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeLabelCmd(cmd parser.Command) error {
	cw.asm.AddComment("label " + cmd.Arg1)
	cw.asm.SetFuncLabel(cw.fnPrefix, cmd.Arg1)
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeIfGotoCmd(cmd parser.Command) error {
	cw.asm.AddComment("if-goto " + cmd.Arg1)
	cw.asm.FromStack("D")
	cw.asm.AtFuncLabel(cw.fnPrefix, cmd.Arg1)
	cw.asm.AsmCmds("D;JNE")
	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeFunctionCmd(cmd parser.Command) error {
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
		cw.asm.AsmCmds(sp, "A=M", "M=0")
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

func (cw *CodeWriter) writeCallCmd(cmd parser.Command) error {
	cw.asm.AddComment(fmt.Sprintf("call %s %d", cmd.Arg1, cmd.Arg2))

	label := fmt.Sprintf("%s.CALL_RET_%d", cw.stPrefix, cw.callCount)
	cw.callCount++

	// Add redturnAddr to stack but do not move SP Pointer
	cw.asm.AtLabel(label)
	cw.asm.AsmCmds("D=A", sp, "A=M", "M=D")
	// Save all segments to the stack except for THAT (the last one)
	segm := [...]segmInstr{lcl, arg, this}
	for _, s := range segm {
		cw.asm.AsmCmds(s, "D=M", sp, "AM=M+1", "M=D")
	}
	// Save THAT to the stack and set SP Pointer to the normal value (empty stack register)
	cw.asm.AsmCmds(that, "D=M", sp, "M=M+1", "M=M+1", "A=M-1", "M=D")
	// Calc new ARG value - it is ARG = SP-5-<func args>
	offset := 5 + cmd.Arg2
	cw.asm.AsmCmds(offset, "D=A", sp, "D=M-D", arg, "M=D")
	// Set new LCL value: LCL=SP
	cw.asm.AsmCmds(sp, "D=M", lcl, "M=D")
	// Jump to called function
	cw.asm.AtLabel(cmd.Arg1)
	cw.asm.AsmCmds("0;JMP")
	// Label of return address
	cw.asm.SetLabel(label)

	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}

func (cw *CodeWriter) writeReturnCmd(cmd parser.Command) error {
	cw.asm.AddComment("return")
	// Save return address. R14 = *(EndFrame - 5)
	cw.asm.AsmCmds(5, "D=A", lcl, "A=M-D", "D=M", "@R14", "M=D")
	// Move return value to arg. *ARG = Pop()
	cw.asm.FromStack("D")
	cw.asm.AsmCmds(arg, "A=M", "M=D")
	// Recycle stack: SP = ARG + 1
	cw.asm.AsmCmds(arg, "D=M+1", sp, "M=D")
	// Restore all func segments from the old stack
	segm := [...]segmInstr{that, this, arg, lcl}
	for _, s := range segm {
		cw.asm.AsmCmds(lcl, "AM=M-1", "D=M", s, "M=D")
	}
	// Jump to return address
	cw.asm.AsmCmds("@R14", "A=M", "0;JMP")

	_, err := cw.writer.WriteString(cw.asm.CodeAsm())
	return err
}
