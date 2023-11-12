package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const INPUT_EX = `
root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32
`

func ParseInput(input string) ([]Monkey, error) {
	var monkeys []Monkey
	scanner := bufio.NewScanner(strings.NewReader(input))
	varPattern := regexp.MustCompile(`^([a-z]+):\s*(\d+)$`)
	opPattern := regexp.MustCompile(`^([a-z]+):\s*([a-z]+)\s*([+*/-])\s*([a-z]+)$`)

	for scanner.Scan() {
		line := scanner.Text()
		if varPattern.MatchString(line) {
			matches := varPattern.FindStringSubmatch(line)
			number, err := strconv.Atoi(matches[2])
			if err != nil {
				return nil, err
			}
			exp := fmt.Sprintf("%d", number)
			if matches[1] == "humn" {
				exp = "x"
			}
			monkeys = append(monkeys, Monkey{
				Name:       matches[1],
				Number:     number,
				LeftDone:   true,
				RightDone:  true,
				Expression: exp,
			})
		} else if opPattern.MatchString(line) {
			matches := opPattern.FindStringSubmatch(line)
			monkeys = append(monkeys, Monkey{
				Name:      matches[1],
				Op:        matches[3],
				LeftName:  matches[2],
				RightName: matches[4],
				LeftDone:  false,
				RightDone: false,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return monkeys, nil
}
