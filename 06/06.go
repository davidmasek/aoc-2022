package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

const SIZE = 14

// Return true if the string contains unique characters.
// Not guarranted to work with unicode.
func isUnique(s string) bool {
	for i, _ := range s {
		for j := i + 1; j < len(s); j++ {
			if s[i] == s[j] {
				return false
			}
		}
	}
	return true
}

func main() {
	fmt.Println(len(os.Args), os.Args)
	fmt.Println("-----")
	if len(os.Args) != 2 {
		fmt.Println("Usage: XY.go XY.in")
		os.Exit(1)
	}
	var filename string
	test := true
	if test {
		filename = strings.Replace(os.Args[1], ".ex", ".in", 1)
	} else {
		filename = os.Args[1]
	}
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		fmt.Println(line)
		for i, _ := range line {
			if isUnique(line[i : i+SIZE]) {
				fmt.Printf("%s at %d\n", line[i:i+SIZE], i+SIZE)
				break
			}
		}
		fmt.Println("-----")
	}
}
