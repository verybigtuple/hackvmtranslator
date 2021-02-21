package main

import "strings"

const (
	// VM commands
	pushKey   = "push"
	popKey    = "pop"
	addKey    = "add"
	subKey    = "sub"
	andKey    = "and"
	orKey     = "or"
	negKey    = "neg"
	notKey    = "not"
	eqKey     = "eq"
	gtKey     = "gt"
	ltKey     = "lt"
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
	addKey: true,
	subKey: true,
	andKey: true,
	orKey:  true,
}

var arithmeticUnaryKeys = map[string]bool{
	negKey: true,
	notKey: true,
}

var arithmeticCondKeys = map[string]bool{
	eqKey: true,
	gtKey: true,
	ltKey: true,
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
