package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var indexColor = color.New(color.FgYellow).SprintFunc()
var occurenceColor = color.New(color.BgYellow, color.FgBlack).SprintFunc()
var fileColor = color.New(color.FgGreen, color.Bold)

var target string
var filesRead int

func main() {
	target = os.Args[1]
	files := Walk(os.Args[2])
	for _, file := range files {
		inspectFile(file)
	}
	fmt.Printf("Read %d files", filesRead)
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
		}
		files = append(files, name)
	}
	return files
}

func inspectFile(path string) {
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
		if http.DetectContentType([]byte(line)) == "application/octet-stream" {
			return
		}
		if strings.Contains(line, target) {
			occurrences = append(occurrences, line)
		}
	}
	if len(occurrences) > 0 {
		fmt.Println()
		fileColor.Println(path)
		fmt.Println()
		for index, i := range occurrences {
			line := strings.Replace(i, target, occurenceColor(target), -1)
			fmt.Printf("%s:%s\n", indexColor(index), line)
		}
	}
}
