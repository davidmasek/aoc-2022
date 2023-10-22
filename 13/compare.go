package main

import (
	"fmt"
	"strings"
)

type Node struct {
	children   []*Node
	isTerminal bool
	value      int
}

type Pair struct {
	left, right *Node
	index       int
}

type Result struct {
	index int
	value string
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

func ReadNode(line string, pos int) (*Node, int) {
	root := NewNode()
	var child *Node
	for line[pos] != ']' {
		if line[pos] == ',' {
			pos++
			continue
		}
		if line[pos] == '[' {
			child, pos = ReadNode(line, pos+1)
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
