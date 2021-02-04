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
	cArithmetic CommandType = iota
	cPush
	cPop

	commentPrefix string = "//"
)

// Command is a struct for a parsed VM cmd
type Command struct {
	CmdType CommandType
	Arg1    string
	Arg2    int
}

func convertTwoArgs(words []string, ct CommandType) (*Command, error) {
	if len(words) < 3 {
		return nil, errors.New("Not enough args for a pop command")
	}
	i, err := strconv.Atoi(words[2])
	if err != nil {
		return nil, err
	}
	return &Command{CmdType: ct, Arg1: words[1], Arg2: i}, nil
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

	switch {
	case ArithmeticKey[words[0]]:
		return &Command{CmdType: cArithmetic, Arg1: words[0]}, nil
	case words[0] == PopKey:
		return convertTwoArgs(words, cPop)
	case words[0] == PushKey:
		return convertTwoArgs(words, cPush)
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
