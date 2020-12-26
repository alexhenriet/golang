package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Parsing parameters
	argv := os.Args
	if len(argv) != 4 {
		fmt.Printf("Error: use %s [original text] [replacement text] [root path]\r\n", argv[0])
		os.Exit(1)
	}
	searchTxt := argv[1]
	replacementTxt := argv[2]
	searchDir := argv[3]
	// Testing parameters
	dir, err := os.Stat(searchDir)
	if err != nil || !dir.IsDir() {
		fmt.Printf("Error: %s must be a readable directory\n", searchDir)
		os.Exit(1)
	}
	// Iterating all files in directory and subdirectories and storing names in fileList
	fileList := collectFiles(searchDir)
	// Using fileList
	for _, file := range fileList {
		if fileContains(file, searchTxt) {
			linesModified := fileReplaceTxt(file, searchTxt, replacementTxt)
			fmt.Printf("%s : %d lines modified\n", file, linesModified)
		}
	}
}

func collectFiles(searchDir string) []string {
	fileList := []string{}
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() { // Skipping directories
			return nil
		}
		fileList = append(fileList, path)
		return nil
	})
	return fileList
}

func fileContains(path string, searchTxt string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), searchTxt) {
			return true
		}
	}
	return false
}

func fileReplaceTxt(path string, searchTxt string, replacementTxt string) int {
	// Opening original file
	f1, err := os.Open(path)
	if err != nil {
		fmt.Printf("%v\n", err)
		return 0
	}
	stats, _ := f1.Stat()
	perms := stats.Mode()
	// Opening modified file
	tmpFile := path + ".tmp"
	f2, err := os.Create(tmpFile)
	if err != nil {
		fmt.Printf("%v\n", err)
		return 0
	}
	f2.Chmod(perms)
	// Writing modified file
	count := 0
	bufr := bufio.NewReader(f1)
	bufr2 := bufio.NewWriter(f2)
	for {
		line, err := bufr.ReadString('\n')
		if strings.Contains(line, searchTxt) {
			line = strings.Replace(line, searchTxt, replacementTxt, -1)
			count++
		}
		bufr2.WriteString(line)
		if err != nil {
			break
		}
	}
	bufr2.Flush()
	f1.Close()
	f2.Close()
	// Erasing original file
	os.Rename(tmpFile, path)
	return count
}
