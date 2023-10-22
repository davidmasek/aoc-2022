package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const LOGGING = false
const TEST = true
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
	var sum int
	fmt.Println("Synchronous")
	fmt.Println("-----")
	start := time.Now()
	results := make([]Result, len(inputs))
	for _, pair := range inputs {
		comp := Compare(pair.left, pair.right)
		results = append(results, Result{pair.index, comp})
	}
	elapsed := time.Since(start)

	debug(results)
	debug("-----")
	sum = 0
	for _, result := range results {
		if result.value == RIGHT {
			sum += result.index
		}
	}
	// A: 5808
	fmt.Println(sum)
	fmt.Println(elapsed)

	fmt.Println("-----")
	fmt.Println("Coroutines")
	fmt.Println("-----")

	// Very naive - creating a goroutine for each input pair
	// Just for fun
	start = time.Now()
	wg := new(sync.WaitGroup)
	resultsChan := make(chan Result, len(inputs))
	for _, pair := range inputs {
		wg.Add(1)
		go func(pair Pair) {
			defer wg.Done()
			comp := Compare(pair.left, pair.right)
			resultsChan <- Result{pair.index, comp}
		}(pair)
	}
	wg.Wait()
	close(resultsChan)
	elapsed = time.Since(start)

	sum = 0
	for result := range resultsChan {
		if result.value == RIGHT {
			sum += result.index
		}
	}
	fmt.Println(sum)
	fmt.Println(elapsed)
}
