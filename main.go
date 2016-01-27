package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var indexColor = color.New(color.FgYellow).SprintFunc()
var occurenceColor = color.New(color.BgYellow, color.FgBlack).SprintFunc()
var fileColor = color.New(color.FgGreen, color.Bold).SprintFunc()

var target string
var filesRead, binaryFiles, matches int

func main() {
	target = os.Args[1]
	files := Walk(os.Args[2])
	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	for index, file := range files {
		if index%1000 == 0 {
			f.Flush()
		}
		inspectFile(file, f)
	}
	f.Write([]byte(fmt.Sprintf("Read %d files", filesRead)))
	f.Write([]byte(fmt.Sprintf("Skipped %d binary files", binaryFiles)))
	f.Write([]byte(fmt.Sprintf("Found %d matches", matches)))
}

func Walk(path string) []string {
	root, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil
	}
	entries, err := root.Readdir(-1)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil
	}

	files := []string{}
	for _, entry := range entries {
		name := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			if strings.HasSuffix(name, "/.git") {
				continue
			}
			files = append(files, Walk(name)...)
			continue
		} else if !entry.Mode().IsRegular() {
			continue
		}
		files = append(files, name)
	}
	return files
}

func inspectFile(path string, f io.Writer) {
	input, err := os.Open(path)
	if err != nil {
		fmt.Printf("File open error: %v", err)
		os.Exit(1)
	}
	occurrences := []string{}
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, target) {
			matches++
			occurrences = append(occurrences, line)
		}
	}
	if len(occurrences) > 0 {
		f.Write([]byte("\n"))
		f.Write([]byte(fmt.Sprintf("%s", fileColor(path))))
		for index, i := range occurrences {
			line := strings.Replace(i, target, occurenceColor(target), -1)
			f.Write([]byte(fmt.Sprintf("%s:%s\n", indexColor(index), line)))
		}
	}
}
