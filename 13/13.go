package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const LOGGING = false
const TEST = false
const B_SIDE = false

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func debug(args ...interface{}) {
	if LOGGING {
		fmt.Println(args...)
	}
}

func main() {
	fmt.Println(len(os.Args), os.Args)
	fmt.Println("-----")
	if len(os.Args) != 2 {
		fmt.Println("Usage: XY.go XY.in")
		os.Exit(1)
	}
	var filename string
	if TEST {
		filename = strings.Replace(os.Args[1], ".ex", ".in", 1)
	} else {
		filename = os.Args[1]
	}
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// reader := bufio.NewReader(f)

	inputs := make([]Pair, 0)
	for i := 1; ; i++ {
		ok := scanner.Scan()
		if !ok {
			break
		}
		line := scanner.Text()
		if line[0] != '[' {
			panic("Expected [")
		}
		left, _ := ReadNode(line, 1)
		ok = scanner.Scan()
		if !ok {
			break
		}
		line = scanner.Text()
		if line[0] != '[' {
			panic("Expected [")
		}
		right, _ := ReadNode(line, 1)

		inputs = append(inputs, Pair{left, right, i})

		ok = scanner.Scan()
		if !ok {
			break
		}
		line = scanner.Text()
		if len(line) != 0 {
			panic("Expected empty line")
		}
	}
	results := make([]Result, 0)
	for _, pair := range inputs {
		debug(pair.left)
		debug(pair.right)
		comp := Compare(pair.left, pair.right)
		debug(comp)
		debug("-----")
		results = append(results, Result{pair.index, comp})
	}

	debug(results)
	debug("-----")
	sum := 0
	for _, result := range results {
		if result.value == RIGHT {
			sum += result.index
		}
	}
	// A: 5808
	fmt.Println(sum)

}
