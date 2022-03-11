package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"regexp"
)

func main() {
	argv := os.Args
	if len(argv) != 2 {
		fmt.Printf("Error: use %s [logsDirectory]\r\n", argv[0])
		os.Exit(1)
	}
	logsDirectory := argv[1]
	dir, err := os.Stat(logsDirectory)
	if err != nil || !dir.IsDir() {
		fmt.Printf("Error: %s must be a readable directory\r\n", logsDirectory)
		os.Exit(1)
	}
    f, err := os.Open(logsDirectory)
    if err != nil {
		fmt.Printf("Error: %v\r\n", err)
		os.Exit(1)
    }
    files, err := f.Readdir(0)
    if err != nil {
		fmt.Printf("Error: %v\r\n", err)
		os.Exit(1)
    }
    for _, currentFile := range files {
		if filepath.Ext(currentFile.Name()) != ".log" {
			continue
		}
		anonymiseLog(logsDirectory, currentFile)
    }
}

func anonymiseLog(directory string, file os.FileInfo) bool {
	logPath := directory + file.Name()
	tmpFilePath := logPath + ".swp"
	// fmt.Printf("Treating %s\r\n", logPath)
	f1, err := os.Open(logPath)
	if err != nil {
		fmt.Printf("Error1: %v\r\n", err)
		os.Exit(1)
	}
	perms := file.Mode()
	stat := file.Sys().(*syscall.Stat_t)
	uid := int(stat.Uid)
	gid := int(stat.Gid)
	f2, err := os.Create(tmpFilePath)
	if err != nil {
		fmt.Printf("Error2: %v\r\n", err)
		os.Exit(1)
	}
	f2.Chmod(perms)
	f2.Chown(uid, gid)

	reader := bufio.NewReader(f1)
	writer := bufio.NewWriter(f2)
	for {
		line, err := reader.ReadString('\n')
		var re = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3})\.\d{1,3}`)
		s := re.ReplaceAllString(line, `$1.0`)
		writer.WriteString(s)
		if err != nil {
			break;
		}
	}
	writer.Flush()
	f1.Close()
	f2.Close()
	os.Rename(tmpFilePath, logPath)
	return true
}