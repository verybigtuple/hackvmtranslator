package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const tempBaseAddr = 5 // Base address for temp vars

// Assign segments to assembler A-Instructions
var segmAInstr = map[string]string{
	"local":    "@LCL",
	"argument": "@ARG",
	"this":     "@THIS",
	"that":     "@THAT",
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

// AddPushFromDReg adds asm code which move SP pointer and push the value of the D-register
// to the stack
func (ah *asmBuilder) AddPushFromDReg() {
	ah.builder.WriteString("@SP\n")
	ah.builder.WriteString("M=M+1\n")
	ah.builder.WriteString("A=M-1\n")
	ah.builder.WriteString("M=D\n")
}

// AddPopToDReg adds asm code which move SP pointer and pop value from the stack to the D-Register
func (ah *asmBuilder) AddPopToDReg() {
	ah.builder.WriteString("@SP\n")
	ah.builder.WriteString("AM=M-1\n")
	ah.builder.WriteString("D=M\n")
}

// AddDeqM adds D=M instruction
func (ah *asmBuilder) AddDeqM() {
	ah.builder.WriteString("D=M\n")
}

// AddMeqD adds M=D instruction
func (ah *asmBuilder) AddMeqD() {
	ah.builder.WriteString("M=D\n")
}

// AddConstToDReg add asm code to add a integer value to the D-Register:
// example
// @101
// D=A
func (ah *asmBuilder) AddConstToDReg(c int) {
	ah.builder.WriteString("@" + strconv.Itoa(c) + "\n")
	ah.builder.WriteString("D=A\n")
}

// AddCalcSegmentAddrWithD adds asm code for calc addr+offset
// by using D register and get addr value to A or D register
//  resReg is A or D (where the calculated addres should be stored)
func (ah *asmBuilder) AddCalcSegmentAddrWithD(segm string, offset int, resReg string) {
	ah.builder.WriteString("@" + strconv.Itoa(offset) + "\n") // like @101
	ah.builder.WriteString("D=A\n")
	ah.builder.WriteString(segmAInstr[segm] + "\n") // like @ARG
	ah.builder.WriteString(resReg + "=D+M\n")
}

// AddCalcSegmentAddr adds asm code that calcs addr+offset w/o using D register
// the result is stored in A register. For big offsets may be ineffective.
func (ah *asmBuilder) AddCalcSegmentAddr(segm string, offset int) {
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

// AddStaticToDReg adds asm code for storing value of a static var to the D register
func (ah *asmBuilder) AddStaticToDReg(prefix string, id int) {
	ah.builder.WriteString(fmt.Sprintf("@%s.%d\n", prefix, id)) //like @file.1
	ah.builder.WriteString("D=M\n")
}

// AddStaticFromDReg adds asm code for moving value of the D-Register to a static var
func (ah *asmBuilder) AddStaticFromDReg(prefix string, id int) {
	ah.builder.WriteString(fmt.Sprintf("@%s.%d\n", prefix, id)) //like @file.1
	ah.builder.WriteString("M=D\n")
}

// AddTempToDReg adds asm code for moving value of a temp var to the D register
func (ah *asmBuilder) AddTempToDReg(offset int) {
	addr := tempBaseAddr + offset
	ah.builder.WriteString("@" + strconv.Itoa(addr) + "\n")
	ah.builder.WriteString("D=M\n")
}

// AddStaticFromDReg adds asm code for moving value of the D-Register to a temp var
func (ah *asmBuilder) AddTempFromDReg(offset int) {
	addr := tempBaseAddr + offset
	ah.builder.WriteString("@" + strconv.Itoa(addr) + "\n")
	ah.builder.WriteString("M=D\n")
}
