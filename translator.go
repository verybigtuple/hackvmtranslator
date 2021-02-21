package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func parseCmdline() (inFilePath string, outFilePath string, err error) {
	inFileFlag := flag.String("in", "", "Input file. Usually has the extension '.vm'")
	outFileFlag := flag.String("out", "", "Output file. Usually has the extension '.asm'")
	flag.Parse()

	inFilePath = *inFileFlag
	if inFilePath == "" {
		if flag.Arg(0) == "" {
			err = fmt.Errorf("Input file is not set")
			return
		}
		inFilePath = flag.Arg(0)
	}

	outFilePath = *outFileFlag
	if outFilePath == "" {
		if flag.Arg(1) == "" {
			fn := strings.TrimSuffix(filepath.Base(inFilePath), filepath.Ext(inFilePath))
			asmFn := fn + ".asm"
			outFilePath = filepath.Join(filepath.Dir(inFilePath), asmFn)
		} else {
			outFilePath = flag.Arg(1)
		}
	}

	return
}

func run(stPrefix string, inReader *bufio.Reader, outWriter *bufio.Writer) error {
	parser := NewParser(inReader)
	codeWr := NewCodeWriter(outWriter, stPrefix, "")
	for {
		cmd, err := parser.ParseNext()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if err := codeWr.WriteCommand(*cmd); err != nil {
			return err
		}
	}
	err := outWriter.Flush()
	return err
}

func main() {
	inFilePath, outFilePath, err := parseCmdline()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Argument Error: %v", err))
		os.Exit(1)
	}

	inFile, err := os.Open(inFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Cannot open input file: %v", err))
		os.Exit(2)
	}
	defer inFile.Close()

	outFile, err := os.Create(outFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Cannot create output file: %v", err))
		os.Exit(2)
	}
	defer func() {
		err := outFile.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Cannot close output file: %v", err))
			os.Exit(2)
		}
	}()

	inReader := bufio.NewReader(inFile)
	outWriter := bufio.NewWriter(outFile)
	stPrefix := strings.TrimSuffix(filepath.Base(inFilePath), filepath.Ext(inFilePath))
	err = run(stPrefix, inReader, outWriter)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Cannot create output file: %v", err))
		os.Exit(3)
	}
	fmt.Printf("Asm file saved as %v\n", outFilePath)
}
