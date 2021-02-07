package main

import (
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
	ah.builder.WriteString(commentPrefix + " " + comment + "\n")
}

func (ah *asmBuilder) AddPushFromDReg() {
	ah.builder.WriteString("@SP\n")
	ah.builder.WriteString("M=M+1\n")
	ah.builder.WriteString("A=M-1\n")
	ah.builder.WriteString("M=D\n")
}

func (ah *asmBuilder) AddPopToDReg() {
	ah.builder.WriteString("@SP\n")
	ah.builder.WriteString("AM=M-1\n")
	ah.builder.WriteString("D=M\n")
}

func (ah *asmBuilder) AddDeqM() {
	ah.builder.WriteString("D=M\n")
}

func (ah *asmBuilder) AddMeqD() {
	ah.builder.WriteString("M=D\n")
}

func (ah *asmBuilder) AddConstToDReg(c int) {
	ah.builder.WriteString("@" + strconv.Itoa(c) + "\n") // like @101
	ah.builder.WriteString("D=A\n")
}

// CalcSegmentAddrWithD returns assemmler code for calc addr+offset
// by using D register and get addr value to A or D register
func (ah *asmBuilder) AddCalcSegmentAddrWithD(segm string, offset int, resReg string) {
	ah.builder.WriteString("@" + strconv.Itoa(offset) + "\n") // like @101
	ah.builder.WriteString("D=A\n")
	ah.builder.WriteString(segmAInstr[segm] + "\n") // like @ARG
	ah.builder.WriteString(resReg + "=D+M\n")
}

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

func (ah *asmBuilder) AddStaticToDReg(prefix string, id int) {
	ah.builder.WriteString(fmt.Sprintf("@%s.%d\n", prefix, id)) //like @file.1
	ah.builder.WriteString("D=M\n")
}

func (ah *asmBuilder) AddStaticFromDReg(prefix string, id int) {
	ah.builder.WriteString(fmt.Sprintf("@%s.%d\n", prefix, id)) //like @file.1
	ah.builder.WriteString("M=D\n")
}

func (ah *asmBuilder) AddTempToDReg(offset int) {
	addr := tempBaseAddr + offset
	ah.builder.WriteString("@" + strconv.Itoa(addr) + "\n")
	ah.builder.WriteString("D=M\n")
}

func (ah *asmBuilder) AddTempFromDReg(offset int) {
	addr := tempBaseAddr + offset
	ah.builder.WriteString("@" + strconv.Itoa(addr) + "\n")
	ah.builder.WriteString("M=D\n")
}
