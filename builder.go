package main

import (
	"bytes"
	"strconv"
	"strings"
)

type SegmInstr string
type freeReg string

const (
	tempBaseAddr = 5 // Base address for temp vars

	SP   SegmInstr = "SP"
	LCL  SegmInstr = "LCL"
	ARG  SegmInstr = "ARG"
	THIS SegmInstr = "THIS"
	THAT SegmInstr = "THAT"

	R13 freeReg = "R13"
	R14 freeReg = "R14"
)

var varSegments = map[string]SegmInstr{
	localKey:    LCL,
	argumentKey: ARG,
	thisKey:     THIS,
	thatKey:     THAT,
}

var pointerOffsets = map[int]SegmInstr{
	0: THIS,
	1: THAT,
}

type asmBuilder struct {
	builder *bytes.Buffer // byte buffer does not reallocates while reset (strings.Builder does)
}

func newAsmBuilder() *asmBuilder {
	b := bytes.Buffer{}
	return &asmBuilder{&b}
}

func (ah *asmBuilder) CodeAsm() string {
	result := ah.builder.String()
	ah.builder.Reset()
	return result
}

func (ah *asmBuilder) AddComment(comment string) {
	if !strings.HasPrefix(comment, commentPrefix) {
		ah.builder.WriteString(commentPrefix + " " + comment + "\n")
	} else {
		ah.builder.WriteString(comment + "\n")
	}
}

// ToStack adds asm code which move SP pointer and push the value of the D-register
// to the stack
func (ah *asmBuilder) ToStack(calc string) {
	ah.AsmCmds(SP, "M=M+1", "A=M-1")
	ah.builder.WriteString("M=")
	ah.builder.WriteString(calc)
	ah.builder.WriteRune('\n')
}

// FromStack adds asm code which move SP pointer and pop value from the stack to the D-Register
func (ah *asmBuilder) FromStack(dest string) {
	ah.AsmCmds(SP, "AM=M-1")
	ah.builder.WriteString(dest)
	ah.builder.WriteString("=M\n")
}

func (ah *asmBuilder) AsmCmds(cmds ...interface{}) {
	for _, c := range cmds {
		switch v := c.(type) {
		case int:
			s := strconv.Itoa(v)
			ah.builder.WriteRune('@')
			ah.builder.WriteString(s)
			ah.builder.WriteRune('\n')
		case freeReg:
			ah.builder.WriteRune('@')
			ah.builder.WriteString(string(v))
			ah.builder.WriteRune('\n')
		case SegmInstr:
			ah.builder.WriteRune('@')
			ah.builder.WriteString(string(v))
			ah.builder.WriteRune('\n')
		case string:
			ah.builder.WriteString(v)
			ah.builder.WriteRune('\n')
		}
	}
}

// StaticAinstr adds a static A-instruction. Like @file.5
func (ah *asmBuilder) StaticAinstr(prefix string, id int) {
	ah.builder.WriteRune('@')
	ah.builder.WriteString(prefix)
	ah.builder.WriteRune('.')
	ah.builder.WriteString(strconv.Itoa(id))
	ah.builder.WriteRune('\n')
}

// StaticAinstr adds a temp A-instruction. Like @7
func (ah *asmBuilder) TempAInstr(offset int) {
	addr := tempBaseAddr + offset
	ah.AsmCmds(addr)
}

// StaticAinstr adds a pointerp A-instruction. @THIS or @THAT
func (ah *asmBuilder) PointerAinstr(offset int) {
	ah.AsmCmds(pointerOffsets[offset])
}

func (ah *asmBuilder) SegmentAinstr(vmSegment string) {
	ah.AsmCmds(varSegments[vmSegment])
}

func (ah *asmBuilder) funcLabel(fnPrefix, label string) {
	ah.builder.WriteString(fnPrefix)
	ah.builder.WriteRune('$')
	ah.builder.WriteString(label)
}

func (ah *asmBuilder) AtFuncLabel(fnPrefix, label string) {
	ah.builder.WriteRune('@')
	ah.funcLabel(fnPrefix, label)
	ah.builder.WriteRune('\n')
}

func (ah *asmBuilder) SetFuncLabel(fnPrefix, label string) {
	ah.builder.WriteRune('(')
	ah.funcLabel(fnPrefix, label)
	ah.builder.WriteString(")\n")
}

func (ah *asmBuilder) arithmCondLabel(statPrefix, cond string, idx int) {
	// static.EQ_END_5
	ah.builder.WriteString(statPrefix)
	ah.builder.WriteRune('.')
	ah.builder.WriteString(strings.ToUpper(cond))
	ah.builder.WriteString("_END_")
	ah.builder.WriteString(strconv.Itoa(idx))
}

func (ah *asmBuilder) AtArithmCondLabel(statPrefix, cond string, idx int) {
	ah.builder.WriteRune('@')
	ah.arithmCondLabel(statPrefix, cond, idx)
	ah.builder.WriteRune('\n')
}

func (ah *asmBuilder) SetArithmCondLabel(statPrefix, cond string, idx int) {
	ah.builder.WriteRune('(')
	ah.arithmCondLabel(statPrefix, cond, idx)
	ah.builder.WriteString(")\n")
}

func (ah *asmBuilder) AtLabel(label string) {
	ah.builder.WriteRune('@')
	ah.builder.WriteString(label)
	ah.builder.WriteRune('\n')
}

func (ah *asmBuilder) SetLabel(label string) {
	ah.builder.WriteRune('(')
	ah.builder.WriteString(label)
	ah.builder.WriteString(")\n")
}
