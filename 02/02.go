package main

import (
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// RULES:
// inputs:
// A, B, C
// X, Y, Z
// rock, paper, scissors
// outcome points:
// loss, draw, win
// 0, 3, 6
const (
	WIN      = 6
	DRAW     = 3
	LOSS     = 0
	ROCK     = "X"
	PAPER    = "Y"
	SCISSORS = "Z"
)

var usagePoints = map[string]int{
	ROCK:     1,
	PAPER:    2,
	SCISSORS: 3,
}
var unify = map[string]string{
	"X": ROCK,
	"Y": PAPER,
	"Z": SCISSORS,
	"A": ROCK,
	"B": PAPER,
	"C": SCISSORS,
}

// Return score for player `b`
func score(a, b string) int {
	a = unify[a]
	b = unify[b]
	if a == b {
		return DRAW
	}
	switch {
	case b == ROCK && a == SCISSORS:
		return WIN
	case b == PAPER && a == ROCK:
		return WIN
	case b == SCISSORS && a == PAPER:
		return WIN
	default:
		return LOSS
	}
}

func main() {
	fmt.Println(len(os.Args), os.Args)
	if len(os.Args) != 2 {
		fmt.Println("Usage: XY.go XY.in")
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	check(err)
	defer f.Close()

	var sum int = 0
	for {
		var opponent, me string
		n, _ := fmt.Fscanln(f, &opponent, &me)
		if n == 2 {
			sum += score(opponent, me) + usagePoints[me]
		} else if n == 0 {
			break
		} else {
			panic("unexpected input format")
		}
	}
	fmt.Println(sum)
}
