package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

type Node struct {
	children   []*Node
	isTerminal bool
	value      int
}

const (
	RIGHT        = "RIGHT"
	UNDETERMINED = "UNDETERMINED"
	WRONG        = "WRONG"
)

// Return if left and right are in the "right" order.
func Compare(left, right *Node) string {
	if left.isTerminal && right.isTerminal {
		if left.value < right.value {
			return RIGHT
		} else if left.value > right.value {
			return WRONG
		} else {
			return UNDETERMINED
		}
	}
	if left.isTerminal {
		left = NewWrapperNode(left.value)

	}
	if right.isTerminal {
		right = NewWrapperNode(right.value)
	}
	i := 0
	for {
		if i >= len(left.children) && i >= len(right.children) {
			return UNDETERMINED
		}
		if i >= len(left.children) {
			return RIGHT
		}
		if i >= len(right.children) {
			return WRONG
		}
		comp := Compare(left.children[i], right.children[i])
		if comp != UNDETERMINED {
			return comp
		}
		i++
	}
}

func (n Node) String() string {
	if n.isTerminal {
		return fmt.Sprintf("%d", n.value)
	} else {
		var children []string
		for _, child := range n.children {
			children = append(children, child.String())
		}
		return fmt.Sprintf("[%s]", strings.Join(children, ","))
	}
}

func NewNode() *Node {
	return &Node{isTerminal: false}
}

func NewWrapperNode(value int) *Node {
	return &Node{isTerminal: false, children: []*Node{NewTerminal(value)}}
}

func NewTerminal(value int) *Node {
	return &Node{isTerminal: true, value: value}
}

func readNode(line string, pos int) (*Node, int) {
	root := NewNode()
	var child *Node
	for line[pos] != ']' {
		if line[pos] == ',' {
			pos++
			continue
		}
		if line[pos] == '[' {
			child, pos = readNode(line, pos+1)
			root.children = append(root.children, child)
		} else {
			var value int
			fmt.Sscanf(line[pos:], "%d", &value)
			root.children = append(root.children, NewTerminal(value))
			pos += len(fmt.Sprint(value))
		}
	}
	return root, pos + 1
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

	results := make([]string, 0)
	for {
		ok := scanner.Scan()
		if !ok {
			break
		}
		line := scanner.Text()
		if line[0] != '[' {
			panic("Expected [")
		}
		left, _ := readNode(line, 1)
		ok = scanner.Scan()
		if !ok {
			break
		}
		line = scanner.Text()
		if line[0] != '[' {
			panic("Expected [")
		}
		right, _ := readNode(line, 1)

		debug(left)
		debug(right)
		comp := Compare(left, right)
		debug(comp)
		debug("-----")
		results = append(results, comp)

		ok = scanner.Scan()
		if !ok {
			break
		}
		line = scanner.Text()
		if len(line) != 0 {
			panic("Expected empty line")
		}
	}
	debug(results)
	debug("-----")
	sum := 0
	for i, result := range results {
		if result == RIGHT {
			sum += i + 1
		}
	}
	// A: 5808
	fmt.Println(sum)

}
