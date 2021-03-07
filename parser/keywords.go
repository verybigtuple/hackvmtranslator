package parser

import "strings"

// VM language KeyWords
const (
	PushKey   = "push"
	PopKey    = "pop"
	AddKey    = "add"
	SubKey    = "sub"
	AndKey    = "and"
	OrKey     = "or"
	NegKey    = "neg"
	NotKey    = "not"
	EqKey     = "eq"
	GtKey     = "gt"
	LtKey     = "lt"
	LabelKey  = "label"
	GotoKey   = "goto"
	IfgotoKey = "if-goto"
	FuncKey   = "function"
	CallKey   = "call"
	ReturnKey = "return"
)

// Segments
const (
	LocalKey    = "local"
	ArgumentKey = "argument"
	ThisKey     = "this"
	ThatKey     = "that"
	ConstantKey = "constant"
	StaticKey   = "static"
	PointerKey  = "pointer"
	TempKey     = "temp"
)

// CommentPrefix is lieteral with that comment starts
const CommentPrefix = "//"

var segmentsKeys = map[string]bool{
	LocalKey:    true,
	ArgumentKey: true,
	ThisKey:     true,
	ThatKey:     true,
	ConstantKey: true,
	StaticKey:   true,
	PointerKey:  true,
	TempKey:     true,
}

func IsConstantSegment(s string) bool {
	return s == ConstantKey
}

func IsStaticSegment(s string) bool {
	return s == StaticKey
}

func IsPointerSegment(s string) bool {
	return s == PointerKey
}

func IsTempSegment(s string) bool {
	return s == TempKey
}

func isValidSegment(s string) bool {
	return segmentsKeys[s]
}

func isValidPopSegment(s string) bool {
	return isValidSegment(s) && !IsConstantSegment(s)
}

func isValidPushSegment(s string) bool {
	return isValidSegment(s)
}

func isComment(s string) bool {
	return strings.HasPrefix(s, CommentPrefix)
}
