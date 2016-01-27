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

var emptyList = []string{}
var target string

func main() {
	target = os.Args[1]
	err := filepath.Walk(os.Args[2], inspectFile)
	if err != nil {
		fmt.Printf("Glob error: %v", err)
		os.Exit(1)
	}
}
func inspectFile(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
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
			return nil
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
	return nil
}
