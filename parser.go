package main

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// CommandType is a type of a parser command
type CommandType int

const (
	// Command Types
	cmdArithmeticBinary CommandType = iota
	cmdArithmeticUnary
	cmdArithmeticCond
	cmdPush
	cmdPop
	cmdGoto
)

// Command is a struct for a parsed VM cmd
type Command struct {
	CmdType CommandType
	Arg1    string
	Arg2    int
}

func convertTwoArgs(words []string, ct CommandType, sValid func(string) bool) (*Command, error) {
	if len(words) < 3 {
		return nil, errors.New("Not enough args for a command")
	}

	if !sValid(words[1]) {
		return nil, fmt.Errorf("Segment argument '%s' is invalid for a command %s", words[1], words[0])
	}

	i, err := strconv.Atoi(words[2])
	if err != nil {
		return nil, fmt.Errorf("Invalid offset argument %s for a command %s: %w", words[2], words[0], err)
	}

	if len(words) > 3 && !isComment(words[3]) {
		return nil, errors.New("Unexpected inline comment literal")
	}

	return &Command{CmdType: ct, Arg1: words[1], Arg2: i}, nil
}

func convertArithmetic(ct CommandType, words []string) (*Command, error) {
	if len(words) > 1 && !isComment(words[1]) {
		return nil, errors.New("Unexpected inline comment literal")
	}
	return &Command{CmdType: ct, Arg1: words[0]}, nil
}

// Parser struct for parsing VM cmds line by line
type Parser struct {
	reader *bufio.Reader
}

func NewParser(r *bufio.Reader) *Parser {
	p := Parser{reader: r}
	return &p
}

func (p Parser) ParseNext() (*Command, error) {
	line, err := p.readNextCodeLine()
	if err != nil {
		return nil, err
	}

	words := strings.Fields(line)
	if len(words) == 0 {
		return nil, fmt.Errorf("No words in the line '%s'", line)
	}

	firstWord := words[0]
	switch {
	case isArithmeticBinary(firstWord):
		return convertArithmetic(cmdArithmeticBinary, words)
	case isArithmeticUnary(firstWord):
		return convertArithmetic(cmdArithmeticUnary, words)
	case isArithmeticCond(firstWord):
		return convertArithmetic(cmdArithmeticCond, words)
	case isPop(firstWord):
		return convertTwoArgs(words, cmdPop, isValidPopSegment)
	case isPush(firstWord):
		return convertTwoArgs(words, cmdPush, isValidPushSegment)
	}

	return nil, fmt.Errorf("Cmd cannot be parsed from line '%s'", line)
}

func (p Parser) readNextLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	// In case the last line does not finish with \n
	if err != nil && len(line) == 0 {
		return "", err
	}
	line = strings.Trim(line, " \t\r\n")
	return line, nil
}

func (p Parser) readNextCodeLine() (string, error) {
	for {
		line, err := p.readNextLine()
		if err != nil {
			return "", err
		}
		if !strings.HasPrefix(line, commentPrefix) && len(line) > 0 {
			return line, nil
		}
	}
}
