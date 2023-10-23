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
	fmt.Println("-----")

	// Part B
	packets := make([]*Node, len(inputs)*2)
	for i, j := 0, 0; i < len(inputs); i, j = i+1, j+2 {
		packets[j] = inputs[i].left
		packets[j+1] = inputs[i].right
	}
	control2, _ := ReadNode("[[2]]", 1)
	packets = append(packets, control2)
	control6, _ := ReadNode("[[6]]", 1)
	packets = append(packets, control6)
	BubbleSort(packets, func(a, b **Node) bool { return Compare(*a, *b) == RIGHT })
	res := 1
	for i, packet := range packets {
		if packet == control2 {
			res *= i + 1
			fmt.Println(i+1, "[[2]]")
		} else if packet == control6 {
			res *= i + 1
			fmt.Println(i+1, "[[6]]")
		}
	}
	debug(packets)
	fmt.Println(res)
	fmt.Println("-----")

	// Part B - alt
	packets = make([]*Node, len(inputs)*2)
	for i, j := 0, 0; i < len(inputs); i, j = i+1, j+2 {
		packets[j] = inputs[i].left
		packets[j+1] = inputs[i].right
	}
	control2, _ = ReadNode("[[2]]", 1)
	control6, _ = ReadNode("[[6]]", 1)
	control2Index := 0
	control6Index := 0
	for _, packet := range packets {
		if Compare(packet, control2) == RIGHT {
			control2Index++
		} else if Compare(packet, control6) == RIGHT {
			control6Index++
		}
	}

	// it's zero based so add 1
	control2Index++
	control6Index++
	// we only checked control6 from those larger than control2
	// so we actually calculated the shift
	control6Index += control2Index
	res = control2Index * control6Index
	fmt.Println(control2Index, "[[2]]")
	fmt.Println(control6Index, "[[6]]")
	fmt.Println(res)
	fmt.Println("-----")
}
