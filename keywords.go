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

	// Segments
	localKey    = "local"
	argumentKey = "argument"
	thisKey     = "this"
	thatKey     = "that"
	constantKey = "constant"
	staticKey   = "static"
	pointerKey  = "pointer"
	tempKey     = "temp"

	commentPrefix = "//"
)

var segmentsKeys = map[string]bool{
	localKey:    true,
	argumentKey: true,
	thisKey:     true,
	thatKey:     true,
	constantKey: true,
	staticKey:   true,
	pointerKey:  true,
	tempKey:     true,
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
