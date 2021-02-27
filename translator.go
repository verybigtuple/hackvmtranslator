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

func getInputFiles(path string) ([]string, error) {
	rootPathInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var matches []string

	if rootPathInfo.IsDir() {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if m, err := filepath.Match("*.vm", path); err != nil {
				return err
			} else if m {
				matches = append(matches, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		matches = append(matches, path)
	}
	return matches, nil
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

func worker(filePath string, result chan<- *strings.Builder, errChan chan<- error) {
	inFile, err := os.Open(filePath)
	if err != nil {
		errChan <- err
	}
	defer inFile.Close()

	fmt.Printf("Reading file %s\n", filePath)

	inReader := bufio.NewReader(inFile)
	sBuilder := &strings.Builder{}
	outWriter := bufio.NewWriter(sBuilder)
	stPrefix := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	err = run(stPrefix, inReader, outWriter)
	if err != nil {
		errChan <- fmt.Errorf("File %s: %w", filePath, err)
		return
	}
	outWriter.Flush()
	result <- sBuilder
}

func main() {
	inFilePath, outFilePath, err := parseCmdline()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Argument Error: %v", err))
		os.Exit(1)
	}

	inPaths, err := getInputFiles(inFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Cannot get input file or directory: %v", err))
		os.Exit(2)
	}

	errChan := make(chan error)
	result := make(chan *strings.Builder)
	for _, inPath := range inPaths {
		go worker(inPath, result, errChan)
	}

	var allErrs []error
	var allRes []*strings.Builder
	for i := 0; i < len(inPaths); i++ {
		select {
		case e := <-errChan:
			allErrs = append(allErrs, e)
		case r := <-result:
			allRes = append(allRes, r)
		}
	}

	if len(allErrs) > 0 {
		fmt.Fprintln(os.Stderr, "Errors during translation:")
		for _, err := range allErrs {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(3)
	}

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

	for _, b := range allRes {
		outFile.WriteString(b.String())
	}
	fmt.Printf("Asm file saved as %v\n", outFilePath)
}
