package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Sprintf("Syntax: %s wordsFile.txt", os.Args[0]))
	}
	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		panic(fmt.Sprintf("%s file does not exist", os.Args[1]))
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Mode [1,2]")
	mode := read(reader)
	if mode != "1" && mode != "2" {
		panic("Invalid mode " + mode)
	}
	m := getMapFromFile(os.Args[1], mode)
	rand.Seed(time.Now().UnixNano())
	for {
		keys := getKeys(&m)
		nbKeys := len(keys)
		if nbKeys == 0 {
			break
		}
		randKey := keys[rand.Intn(nbKeys)]
		for {
			fmt.Printf("%s -> ", randKey)
			answer := read(reader)
			if strings.EqualFold(answer, m[randKey]) {
				break
			} else if answer == "?" {
				fmt.Printf("%s\n", m[randKey])
			}
		}
		delete(m, randKey)
	}
}

func read(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}

func getMapFromFile(filename string, mode string) map[string]string {
	var m = make(map[string]string)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), ":")
		left := strings.TrimSpace(split[0])
		right := strings.TrimSpace(split[1])
		if mode == "1" {
			m[left] = right
		} else {
			m[right] = left
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return m
}

func getKeys(m *map[string]string) []string {
	var keys []string
	for k := range *m {
		keys = append(keys, k)
	}
	return keys
}
