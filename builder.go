package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const (
	tempBaseAddr = 5 // Base address for temp vars
	freeReg      = "@R13"
)

// Assign segments to assembler A-Instructions
var segmAInstr = map[string]string{
	"local":    "@LCL",
	"argument": "@ARG",
	"this":     "@THIS",
	"that":     "@THAT",
}

var pointerOffsets = map[int]string{
	0: "@THIS",
	1: "@THAT",
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

// FromDtoStack adds asm code which move SP pointer and push the value of the D-register
// to the stack
func (ah *asmBuilder) FromDtoStack() {
	ah.builder.WriteString("@SP\n")
	ah.builder.WriteString("M=M+1\n")
	ah.builder.WriteString("A=M-1\n")
	ah.builder.WriteString("M=D\n")
}

// FromStackToD adds asm code which move SP pointer and pop value from the stack to the D-Register
func (ah *asmBuilder) FromStackToD() {
	ah.builder.WriteString("@SP\n")
	ah.builder.WriteString("AM=M-1\n")
	ah.builder.WriteString("D=M\n")
}

func (ah *asmBuilder) SetTopStack(calc string) {
	ah.builder.WriteString("@SP\n")
	ah.builder.WriteString("A=M-1\n")
	ah.builder.WriteString("M=" + calc + "\n")
}

func (ah *asmBuilder) DecAddr() {
	ah.builder.WriteString("A=A-1\n")
}

func (ah *asmBuilder) ArbitraryCmd(cmd string) {
	ah.builder.WriteString(cmd + "\n")
}

func (ah *asmBuilder) CondFalseDefault() {
	ah.builder.WriteString("A=A-1\n")
	ah.builder.WriteString("D=M-D\n")
	ah.builder.WriteString("M=0\n")
}

func (ah *asmBuilder) CondJump(prefix, cond string, c int, jmp string) {
	up := strings.ToUpper(cond)
	label := fmt.Sprintf("%s.%s_END_%d", prefix, up, c)
	ah.builder.WriteString("@" + label + "\n")
	ah.builder.WriteString("D;" + jmp + "\n")
	ah.builder.WriteString("@SP\n")
	ah.builder.WriteString("A=M-1\n")
	ah.builder.WriteString("M=-1\n")
	ah.builder.WriteString("(" + label + ")\n")
}

// FromMemToD adds D=M instruction
func (ah *asmBuilder) FromMemToD() {
	ah.builder.WriteString("D=M\n")
}

// FromDtoMem adds M=D instruction
func (ah *asmBuilder) FromDtoMem() {
	ah.builder.WriteString("M=D\n")
}

// AddRregFromDReg - asm code for adding a value of the D-Register to R-register
func (ah *asmBuilder) ToR(from string) {
	ah.builder.WriteString(freeReg + "\n")
	ah.builder.WriteString("M=" + from + "\n")
}

func (ah *asmBuilder) FromR(src string) {
	ah.builder.WriteString(freeReg + "\n")
	ah.builder.WriteString(src + "=M\n")
}

// ConstToD add asm code to add a integer value to the D-Register:
// example
// @101
// D=A
func (ah *asmBuilder) ConstToD(c int) {
	ah.builder.WriteString("@" + strconv.Itoa(c) + "\n")
	ah.builder.WriteString("D=A\n")
}

// SegmAddrCalcWithD adds asm code for calc addr+offset
// by using D register and get addr value to A or D register
//  resReg is A or D (where the calculated addres should be stored)
func (ah *asmBuilder) SegmAddrCalcWithD(segm string, offset int, resReg string) {
	ah.builder.WriteString("@" + strconv.Itoa(offset) + "\n") // like @101
	ah.builder.WriteString("D=A\n")
	ah.builder.WriteString(segmAInstr[segm] + "\n") // like @ARG
	ah.builder.WriteString(resReg + "=D+M\n")
}

// SegmAddr adds asm code that calcs addr+offset w/o using D register
// the result is stored in A register. For big offsets may be ineffective.
func (ah *asmBuilder) SegmAddr(segm string, offset int) {
	ah.builder.WriteString(segmAInstr[segm] + "\n")
	if offset == 0 {
		ah.builder.WriteString("A=M\n")
	} else {
		ah.builder.WriteString("A=M+1\n")
	}
	for i := 0; i < offset-1; i++ {
		ah.builder.WriteString("A=A+1\n")
	}
}

// StaticToD adds asm code for storing value of a static var to the D register
func (ah *asmBuilder) StaticToD(prefix string, id int) {
	ah.builder.WriteString(fmt.Sprintf("@%s.%d\n", prefix, id)) //like @file.1
	ah.builder.WriteString("D=M\n")
}

// StaticFromD adds asm code for moving value of the D-Register to a static var
func (ah *asmBuilder) StaticFromD(prefix string, id int) {
	ah.builder.WriteString(fmt.Sprintf("@%s.%d\n", prefix, id)) //like @file.1
	ah.builder.WriteString("M=D\n")
}

// TempToD adds asm code for moving value of a temp var to the D register
func (ah *asmBuilder) TempToD(offset int) {
	addr := tempBaseAddr + offset
	ah.builder.WriteString("@" + strconv.Itoa(addr) + "\n")
	ah.builder.WriteString("D=M\n")
}

// AddStaticFromDReg adds asm code for moving value of the D-Register to a temp var
func (ah *asmBuilder) TempFromD(offset int) {
	addr := tempBaseAddr + offset
	ah.builder.WriteString("@" + strconv.Itoa(addr) + "\n")
	ah.builder.WriteString("M=D\n")
}

func (ah *asmBuilder) PointerToD(offset int) {
	ah.builder.WriteString(pointerOffsets[offset] + "\n")
	ah.builder.WriteString("D=M\n")
}

func (ah *asmBuilder) PointerFromD(offset int) {
	ah.builder.WriteString(pointerOffsets[offset] + "\n")
	ah.builder.WriteString("M=D\n")
}

func (ah *asmBuilder) AtLabel(fnPrefix, label string) {
	if fnPrefix == "" {
		ah.builder.WriteString("@" + label + "\n")
	} else {
		ah.builder.WriteString("@" + fnPrefix + "$" + label + "\n")
	}
}

func (ah *asmBuilder) SetLabel(fnPrefix, label string) {
	ah.builder.WriteString("(" + fnPrefix + "$" + label + ")\n")
}
