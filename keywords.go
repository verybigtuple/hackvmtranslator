package main

import "strings"

const (
	pushKey = "push"
	popKey  = "pop"

	constantKey = "constant"
	staticKey   = "static"
	pointerKey  = "pointer"
	tempKey     = "temp"

	commentPrefix = "//"
)

var arithmeticKeys = map[string]bool{
	"add": true,
	"sub": true,
	"neg": true,
	"eq":  true,
	"gt":  true,
	"lt":  true,
	"and": true,
	"or":  true,
	"not": true,
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
	return arithmeticKeys[s]
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
