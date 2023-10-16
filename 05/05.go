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

type Queue struct {
	elements []rune
}

func (s *Queue) PushRight(e rune) {
	s.elements = append(s.elements, e)
}

func (s *Queue) PushLeft(e rune) {
	s.elements = append([]rune{e}, s.elements...)
}

func (s *Queue) PopRight() rune {
	if len(s.elements) == 0 {
		panic("Stack is empty")
	}
	e := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return e
}

func (s *Queue) PopLeft() rune {
	if len(s.elements) == 0 {
		panic("Stack is empty")
	}
	e := s.elements[0]
	s.elements = s.elements[1:]
	return e
}

func (s *Queue) Content() string {
	return string(s.elements)
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

	// we're limited by design to single digit numbers
	var cols [10]Queue

	// CrateMover 9001
	b_side := true
	var aux Queue

	var n_cols int
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		// each col takes 4 characters including space
		// there is space missing at the end
		n_cols = (len(line) + 1) / 4
		for c := 0; c < n_cols; c++ {
			i := c*4 + 1
			char := rune(line[i])
			if char != ' ' {
				cols[c].PushLeft(char)
			}
		}
	}
	// we always read one extra entry because of column labels
	for c := 0; c < n_cols; c++ {
		cols[c].PopLeft()
	}
	for c := 0; c < n_cols; c++ {
		fmt.Println(c, cols[c].Content())
	}
	fmt.Println("-----")
	for scanner.Scan() {
		line := scanner.Text()
		var from, to, count int
		_, err = fmt.Sscanf(line, "move %d from %d to %d", &count, &from, &to)
		check(err)
		// we're using 0-based indexing
		from--
		to--
		for i := 0; i < count; i++ {
			e := cols[from].PopRight()
			if b_side {
				aux.PushRight(e)
			} else {
				cols[to].PushRight(e)
			}
		}
		if b_side {
			for i := 0; i < count; i++ {
				e := aux.PopRight()
				cols[to].PushRight(e)
			}
		}

		// fmt.Println("move", count, "from", from, "to", to)
		// for c := 0; c < n_cols; c++ {
		// 	fmt.Println(c, cols[c].Content())
		// }
		// fmt.Println("-----")
	}
	for c := 0; c < n_cols; c++ {
		fmt.Println(c, cols[c].Content())
	}
	fmt.Println("-----")
	for c := 0; c < n_cols; c++ {
		fmt.Print(string(cols[c].PopRight()))
	}
	fmt.Println()
}
