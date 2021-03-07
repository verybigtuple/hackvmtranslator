package main

import (
	"bufio"
	"container/heap"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func parseCmdline() (inPath, outFilePath string, noBootstrap bool, err error) {
	inFileFlag := flag.String("in", "", "Input file or folder with *.vm files")
	outFileFlag := flag.String("out", "", "Output file. Usually has the extension '.asm'")
	flag.BoolVar(
		&noBootstrap,
		"nb",
		false,
		"Translator does not write the bootstrapping code to a result asm file",
	)
	flag.Parse()

	inPath = *inFileFlag
	if inPath == "" {
		if flag.Arg(0) == "" {
			err = fmt.Errorf("Input file/folder is not set")
			return
		}
		inPath = flag.Arg(0)
	}

	outFilePath = *outFileFlag
	if outFilePath == "" {
		if flag.Arg(1) == "" {
			info, erri := os.Stat(inPath)
			if erri != nil {
				err = fmt.Errorf("Illegal input path: %w", erri)
				return
			}
			if info.IsDir() {
				outFilePath = filepath.Join(inPath, filepath.Base(inPath)+".asm")
			} else {
				fn := strings.TrimSuffix(filepath.Base(inPath), filepath.Ext(inPath))
				asmFn := fn + ".asm"
				outFilePath = filepath.Join(filepath.Dir(inPath), asmFn)
			}
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
			if m, err := filepath.Match("*.vm", filepath.Base(path)); err != nil {
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

func run(writerName, stPrefix string, inReader *bufio.Reader, outWriter *bufio.Writer) error {
	parser := NewParser(inReader)
	codeWr := NewCodeWriter(outWriter, writerName, stPrefix, "")
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

func processBootstrap(result chan<- *trResult, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	sBuilder := &strings.Builder{}
	outWriter := bufio.NewWriter(sBuilder)
	bsCodeWriter := NewCodeWriterBootstrap(outWriter)
	err := bsCodeWriter.WriteBootstrap()
	if err != nil {
		errChan <- err
		return
	}
	outWriter.Flush()
	result <- &trResult{bootstrap, sBuilder}
}

func processVMFile(
	filePath string,
	result chan<- *trResult,
	errChan chan<- error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	inFile, err := os.Open(filePath)
	if err != nil {
		errChan <- err
	}
	defer inFile.Close()

	fmt.Printf("Reading file %s\n", filePath)

	inReader := bufio.NewReader(inFile)
	sBuilder := &strings.Builder{}
	outWriter := bufio.NewWriter(sBuilder)

	fBase := filepath.Base(filePath)
	stPrefix := strings.TrimSuffix(fBase, filepath.Ext(filePath))
	err = run(fBase, stPrefix, inReader, outWriter)
	if err != nil {
		errChan <- fmt.Errorf("File %s: %w", filePath, err)
		return
	}
	outWriter.Flush()
	result <- &trResult{Name: fBase, Builder: sBuilder}
}

func gatherResults(r <-chan *trResult, e <-chan error, wg *sync.WaitGroup) (*resPriotityQueue, []error) {
	es := []error{}
	rq := resPriotityQueue{}
	heap.Init(&rq)

	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

RLoop:
	for {
		select {
		case err := <-e:
			es = append(es, err)
		case res := <-r:
			heap.Push(&rq, res)
		case <-done:
			break RLoop
		}
	}
	return &rq, es
}

func writeAsmFile(filePath string, rq *resPriotityQueue) (err error) {
	outFile, err := os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("Cannot create output file: %w", err)
		return
	}
	defer func() {
		cerr := outFile.Close()
		if cerr != nil {
			err = fmt.Errorf("Cannot close output file: %w", cerr)
		}
	}()

	for len(*rq) > 0 {
		r := heap.Pop(rq).(*trResult)
		_, err = outFile.WriteString(r.Builder.String())
	}
	fmt.Printf("Asm file saved as %v\n", filePath)
	return
}

func main() {
	inFilePath, outFilePath, noBoot, err := parseCmdline()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Argument Error: %v", err))
		os.Exit(1)
	}

	inPaths, err := getInputFiles(inFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Cannot get input file or directory: %v", err))
		os.Exit(2)
	}

	resChan := make(chan *trResult)
	errChan := make(chan error)
	wg := &sync.WaitGroup{}

	if !noBoot {
		wg.Add(1)
		go processBootstrap(resChan, errChan, wg)
	}
	for _, inPath := range inPaths {
		wg.Add(1)
		go processVMFile(inPath, resChan, errChan, wg)
	}

	resultQueue, allErrs := gatherResults(resChan, errChan, wg)
	if len(allErrs) > 0 {
		fmt.Fprintln(os.Stderr, "Errors during translation:")
		for _, err := range allErrs {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(3)
	}
	err = writeAsmFile(outFilePath, resultQueue)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}
