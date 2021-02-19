package main

import "strings"

const (
	pushKey = "push"
	popKey  = "pop"

	labelKey  = "label"
	gotoKey   = "goto"
	ifgotoKey = "if-goto"

	constantKey = "constant"
	staticKey   = "static"
	pointerKey  = "pointer"
	tempKey     = "temp"

	commentPrefix = "//"
)

var arithmeticBinaryKeys = map[string]bool{
	"add": true,
	"sub": true,
	"and": true,
	"or":  true,
}

var arithmeticUnaryKeys = map[string]bool{
	"neg": true,
	"not": true,
}

var arithmeticCondKeys = map[string]bool{
	"eq": true,
	"gt": true,
	"lt": true,
}

var segmentsKeys = map[string]bool{
	"local":     true,
	"argument":  true,
	"this":      true,
	"that":      true,
	constantKey: true,
	staticKey:   true,
	pointerKey:  true,
	tempKey:     true,
}

func isPush(s string) bool {
	return s == pushKey
}

func isPop(s string) bool {
	return s == popKey
}

func isArithmetic(s string) bool {
	return arithmeticBinaryKeys[s] || arithmeticUnaryKeys[s] || arithmeticCondKeys[s]
}

func isArithmeticBinary(s string) bool {
	return arithmeticBinaryKeys[s]
}

func isArithmeticUnary(s string) bool {
	return arithmeticUnaryKeys[s]
}

func isArithmeticCond(s string) bool {
	return arithmeticCondKeys[s]
}

func isGoto(s string) bool {
	return s == gotoKey
}

func isLabel(s string) bool {
	return s == labelKey
}

func isIfGoto(s string) bool {
	return s == ifgotoKey
}

func isConstantSegment(s string) bool {
	return s == constantKey
}

func isStaticSegment(s string) bool {
	return s == staticKey
}

func isPointerSegment(s string) bool {
	return s == pointerKey
}

func isTempSegment(s string) bool {
	return s == tempKey
}

func isValidSegment(s string) bool {
	return segmentsKeys[s]
}

func isValidPopSegment(s string) bool {
	return isValidSegment(s) && !isConstantSegment(s)
}

func isValidPushSegment(s string) bool {
	return isValidSegment(s)
}

func isComment(s string) bool {
	return strings.HasPrefix(s, commentPrefix)
}
