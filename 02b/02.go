package main

import (
	"fmt"
	"os"
	"strings"
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
	WIN_SCORE  = 6
	DRAW_SCORE = 3
	LOSS_SCORE = 0
	LOSS       = "X"
	DRAW       = "Y"
	WIN        = "Z"
	ROCK       = "A"
	PAPER      = "B"
	SCISSORS   = "C"
)

var usagePoints = map[string]int{
	ROCK:     1,
	PAPER:    2,
	SCISSORS: 3,
}
var useToWin = map[string]string{
	ROCK:     PAPER,
	PAPER:    SCISSORS,
	SCISSORS: ROCK,
}
var useToLose = map[string]string{
	ROCK:     SCISSORS,
	PAPER:    ROCK,
	SCISSORS: PAPER,
}

func score(opponent, result string) int {
	switch result {
	case DRAW:
		return DRAW_SCORE + usagePoints[opponent]
	case WIN:
		return WIN_SCORE + usagePoints[useToWin[opponent]]
	case LOSS:
		return LOSS_SCORE + usagePoints[useToLose[opponent]]
	default:
		panic("unexpected result")
	}
}

func main() {
	fmt.Println(len(os.Args), os.Args)
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

	var sum int = 0
	for {
		var opponent, result string
		n, _ := fmt.Fscanln(f, &opponent, &result)
		if n == 2 {
			sum += score(opponent, result)
		} else if n == 0 {
			break
		} else {
			panic("unexpected input format")
		}
	}
	fmt.Println(sum)
}
