package main

import (
	"fmt"
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
				monkeys[name].LeftName = ""
				monkeys[name].LeftValue = leaf.Number
			} else if leaf.Name == monkeys[name].RightName {
				monkeys[name].RightName = ""
				monkeys[name].RightValue = leaf.Number
			} else {
				panic("Incorrect state")
			}
			if monkeys[name].LeftName == "" && monkeys[name].RightName == "" {
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
				leafs = append(leafs, monkeys[name])
			}
		}
	}
	fmt.Println("# Opened:", counter)
	fmt.Println(monkeys["root"])
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
