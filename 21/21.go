package main

import (
	"fmt"
	"strings"
)

type Monkey struct {
	Name       string
	Number     int
	Op         string
	LeftName   string
	RightName  string
	Blocks     []string
	LeftValue  int
	RightValue int
	LeftDone   bool
	RightDone  bool
	Expression string
}

func solve(monkeysArr []Monkey) {
	fmt.Println("# Input:", len(monkeysArr))
	monkeys := make(map[string]*Monkey)
	for i := range monkeysArr {
		m := &monkeysArr[i]
		monkeys[monkeysArr[i].Name] = m
		m.Blocks = make([]string, 0)
	}
	for _, m := range monkeys {
		if m.LeftName != "" {
			monkeys[m.LeftName].Blocks = append(monkeys[m.LeftName].Blocks, m.Name)
		}
		if m.RightName != "" {
			monkeys[m.RightName].Blocks = append(monkeys[m.RightName].Blocks, m.Name)
		}
	}

	leafs := make([]*Monkey, 0)
	for _, m := range monkeys {
		if m.LeftName == "" && m.RightName == "" {
			leafs = append(leafs, m)
		}
	}
	counter := 0
	for len(leafs) > 0 {
		counter++
		leaf := leafs[0]
		leafs = leafs[1:]
		for _, name := range leaf.Blocks {
			if leaf.Name == monkeys[name].LeftName {
				monkeys[name].LeftDone = true
				monkeys[name].LeftValue = leaf.Number
			} else if leaf.Name == monkeys[name].RightName {
				monkeys[name].RightDone = true
				monkeys[name].RightValue = leaf.Number
			} else {
				panic("Incorrect state")
			}
			if monkeys[name].LeftDone && monkeys[name].RightDone {
				switch monkeys[name].Op {
				case "+":
					monkeys[name].Number = monkeys[name].LeftValue + monkeys[name].RightValue
				case "-":
					monkeys[name].Number = monkeys[name].LeftValue - monkeys[name].RightValue
				case "*":
					monkeys[name].Number = monkeys[name].LeftValue * monkeys[name].RightValue
				case "/":
					monkeys[name].Number = monkeys[name].LeftValue / monkeys[name].RightValue
				default:
					panic(leaf.Op)

				}
				current := monkeys[name]
				if current.Name == "root" {
					monkeys[name].Expression = fmt.Sprintf("%s\n================\n%s",
						monkeys[current.LeftName].Expression,
						monkeys[current.RightName].Expression,
					)
					fmt.Println("[root] Expression:", monkeys[name].Expression)
				} else if current.Name == "humn" {
					// actually set during loading, but this should be more general in case of other inputs
					monkeys[name].Expression = "x"
					fmt.Println("[humn] Expression:", monkeys[name].Expression)
				} else {
					left := monkeys[current.LeftName].Expression
					if !strings.Contains(left, "x") {
						left = fmt.Sprintf("%d", monkeys[current.LeftName].Number)
					}
					right := monkeys[current.RightName].Expression
					if !strings.Contains(right, "x") {
						right = fmt.Sprintf("%d", monkeys[current.RightName].Number)
					}
					monkeys[name].Expression = fmt.Sprintf("(%s %s %s)",
						left,
						monkeys[name].Op,
						right,
					)
				}
				leafs = append(leafs, monkeys[name])
			}
		}
	}
	fmt.Println("# Opened:", counter)
	fmt.Println("A:", monkeys["root"].Number)
	// log.Fatal("Done")
}

func main() {
	monkeys, err := ParseInput(INPUT_EX)
	if err != nil {
		panic(err)
	}
	solve(monkeys)

	monkeys, err = ParseInput(INPUT_TEST)
	if err != nil {
		panic(err)
	}
	solve(monkeys)
}
