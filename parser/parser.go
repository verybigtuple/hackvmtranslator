package parser

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// CommandType is a type of a parser command
type CommandType int

// All command types of VM language
const (
	CmdArithmeticBinary CommandType = iota
	CmdArithmeticUnary
	CmdArithmeticCond
	CmdPush
	CmdPop
	CmdLabel
	CmdGoto
	CmdIfGoto
	CmdFunction
	CmdCall
	CmdReturn
)

var cmdTypes = map[string]CommandType{
	PushKey:   CmdPush,
	PopKey:    CmdPop,
	AddKey:    CmdArithmeticBinary,
	SubKey:    CmdArithmeticBinary,
	AndKey:    CmdArithmeticBinary,
	OrKey:     CmdArithmeticBinary,
	NegKey:    CmdArithmeticUnary,
	NotKey:    CmdArithmeticUnary,
	EqKey:     CmdArithmeticCond,
	GtKey:     CmdArithmeticCond,
	LtKey:     CmdArithmeticCond,
	LabelKey:  CmdLabel,
	GotoKey:   CmdGoto,
	IfgotoKey: CmdIfGoto,
	FuncKey:   CmdFunction,
	CallKey:   CmdCall,
	ReturnKey: CmdReturn,
}

var cmdConverters = map[CommandType]func(CommandType, []string) (*Command, error){
	CmdPush:             convertPushPop,
	CmdPop:              convertPushPop,
	CmdArithmeticBinary: convertArithmetic,
	CmdArithmeticUnary:  convertArithmetic,
	CmdArithmeticCond:   convertArithmetic,
	CmdLabel:            conevrtLabeled,
	CmdGoto:             conevrtLabeled,
	CmdIfGoto:           conevrtLabeled,
	CmdFunction:         convertFunc,
	CmdCall:             convertFunc,
	CmdReturn:           convertReturn,
}

func checkNullArgs(words []string) (err error) {
	if len(words) > 1 && !isComment(words[1]) {
		err = errors.New("Too many arguments")
	}
	return
}

func checkOneArg(words []string) (string, error) {
	if len(words) < 2 || (len(words) > 2 && !isComment(words[2])) {
		return "", errors.New("One argument expected")
	}
	return words[1], nil
}

func checkTwoArgs(words []string) (string, int, error) {
	if len(words) < 3 || (len(words) > 3 && !isComment(words[3])) {
		return "", 0, errors.New("Three argument expected")
	}
	i, err := strconv.Atoi(words[2])
	if err != nil {
		return "", 0, fmt.Errorf("Second argument %s is not an integer number", words[2])
	}
	return words[1], i, nil
}

func convertPushPop(ct CommandType, words []string) (*Command, error) {
	segment, offset, err := checkTwoArgs(words)
	if err != nil {
		return nil, err
	}
	if !(ct == CmdPush && isValidPushSegment(segment)) && !isValidPopSegment(segment) {
		return nil, fmt.Errorf("Invalid segment %s for %s command", segment, words[0])
	}
	if offset < 0 {
		return nil, fmt.Errorf("Offset cannot be negative")
	}
	return &Command{ct, segment, offset}, nil
}

func convertArithmetic(ct CommandType, words []string) (*Command, error) {
	err := checkNullArgs(words)
	if err != nil {
		return nil, err
	}
	return &Command{CmdType: ct, Arg1: words[0]}, nil
}

func conevrtLabeled(ct CommandType, words []string) (*Command, error) {
	label, err := checkOneArg(words)
	if err != nil {
		return nil, err
	}
	return &Command{CmdType: ct, Arg1: label}, nil
}

func convertFunc(ct CommandType, words []string) (*Command, error) {
	funcName, offset, err := checkTwoArgs(words)
	if err != nil {
		return nil, err
	}
	if offset < 0 {
		return nil, fmt.Errorf("Offset cannot be negative")
	}
	return &Command{ct, funcName, offset}, nil
}

func convertReturn(ct CommandType, words []string) (*Command, error) {
	err := checkNullArgs(words)
	if err != nil {
		return nil, err
	}
	return &Command{CmdType: ct}, nil
}

// Command is a struct for a parsed VM cmd
type Command struct {
	CmdType CommandType
	Arg1    string
	Arg2    int
}

// Parser struct for parsing VM cmds line by line
type Parser struct {
	reader *bufio.Reader
	lCount int
}

func NewParser(r *bufio.Reader) *Parser {
	p := Parser{reader: r}
	return &p
}

func (p *Parser) ParseNext() (*Command, error) {
	line, err := p.readNextCodeLine()
	if err != nil {
		return nil, err
	}

	words := strings.Fields(line)
	if len(words) == 0 {
		return nil, fmt.Errorf("Line %d: No words parsed", p.lCount)
	}

	firstWord := words[0]
	cmdType, ok := cmdTypes[firstWord]
	if !ok {
		return nil, fmt.Errorf("Line %d: Unknown command %s", p.lCount, firstWord)
	}

	converter, ok := cmdConverters[cmdType]
	if !ok {
		return nil, fmt.Errorf(
			"Line %d: Cannot parse line '%s' as converter is not set",
			p.lCount, line,
		)
	}
	cmd, err := converter(cmdType, words)
	if err != nil {
		return nil, fmt.Errorf("Line %d: %w", p.lCount, err)
	}

	return cmd, nil
}

func (p *Parser) readNextLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	// In case the last line does not finish with \n
	if err != nil && len(line) == 0 {
		return "", err
	}
	p.lCount++
	line = strings.Trim(line, " \t\r\n")
	return line, nil
}

func (p *Parser) readNextCodeLine() (string, error) {
	for {
		line, err := p.readNextLine()
		if err != nil {
			return "", err
		}
		if !strings.HasPrefix(line, CommentPrefix) && len(line) > 0 {
			return line, nil
		}
	}
}
