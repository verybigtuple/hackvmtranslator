package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		panic("Not enough args")
	}

	inFilePath := os.Args[1]
	inFile, err := os.Open(inFilePath)
	if err != nil {
		panic("No File found")
	}
	defer inFile.Close()
	inReader := bufio.NewReader(inFile)

	outFile, err := os.Create(os.Args[2])
	if err != nil {
		panic("Cannot create output file")
	}
	defer func() {
		err := outFile.Close()
		if err != nil {
			panic("Cannot close output file")
		}
	}()
	outWriter := bufio.NewWriter(outFile)

	stPrefix := strings.TrimSuffix(filepath.Base(inFilePath), filepath.Ext(inFilePath))
	parser := NewParser(inReader)
	codeWr := NewCodeWriter(stPrefix, outWriter)
	for {
		cmd, err := parser.ParseNext()
		// println(cmd)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		if err := codeWr.WriteCommand(*cmd); err != nil {
			panic(err)
		}
	}
	outWriter.Flush()
	outFile.Sync()
}
