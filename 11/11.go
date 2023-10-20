package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

const LOGGING = false

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Monkey struct {
	id          int
	items       []int
	op          string
	factor      int
	testval     int
	trueTarget  int
	falseTarget int
	activity    int
}

// Read one monkey. Example input below.
//
//	 Monkey 0:
//		Starting items: 79, 98
//		Operation: new = old * 19
//		Test: divisible by 23
//		  If true: throw to monkey 2
//		  If false: throw to monkey 3
func readMonkey(scanner *bufio.Scanner) (Monkey, bool) {
	// --- monkey id
	scanner.Scan()
	line := scanner.Text()
	var id int
	fmt.Sscanf(line, "Monkey %d:", &id)
	fmt.Println("Monkey", id)

	// --- starting items
	var items []int
	scanner.Scan()
	line = scanner.Text()
	lineItems := strings.Split(line, " ")

	for i := range lineItems {
		if len(lineItems[i]) > 0 && unicode.IsNumber(rune(lineItems[i][0])) {
			var item int
			fmt.Sscanf(lineItems[i], "%d,", &item)
			items = append(items, item)
		}
	}
	fmt.Println("Items", items)
	// --- operation
	var op, old string
	var factor int
	scanner.Scan()
	line = scanner.Text()
	fmt.Sscanf(line, "  Operation: new = old %s %s", &op, &old)
	if old == "old" {
		factor = -1
	} else {
		fmt.Sscanf(line, "  Operation: new = old %s %d", &op, &factor)
	}
	fmt.Println("Operation", op, factor)

	// --- test
	var testval int
	scanner.Scan()
	line = scanner.Text()
	fmt.Sscanf(line, "  Test: divisible by %d", &testval)
	fmt.Println("Testval", testval)

	// --- true/false targets
	var trueTarget, falseTarget int
	scanner.Scan()
	line = scanner.Text()
	fmt.Sscanf(line, "    If true: throw to monkey %d", &trueTarget)
	scanner.Scan()
	line = scanner.Text()
	fmt.Sscanf(line, "    If false: throw to monkey %d", &falseTarget)
	fmt.Println("Targets", trueTarget, falseTarget)

	m := Monkey{id, items, op, factor, testval, trueTarget, falseTarget, 0}
	// empty line between monkeys
	ok := scanner.Scan()
	_ = scanner.Text()
	return m, ok
}

func debug(args ...interface{}) {
	if LOGGING {
		fmt.Println(args...)
	}
}

func max(arr []int) (int, int) {
	if len(arr) == 0 {
		panic("empty array")
	}
	maxIndex := 0
	max := arr[0]
	for i := range arr {
		if arr[i] > max {
			max = arr[i]
			maxIndex = i
		}
	}
	return maxIndex, max
}

func main() {
	fmt.Println(len(os.Args), os.Args)
	fmt.Println("-----")
	if len(os.Args) != 2 {
		fmt.Println("Usage: XY.go XY.in")
		os.Exit(1)
	}
	var filename string
	const TEST = true
	const B_SIDE = true
	if TEST {
		filename = strings.Replace(os.Args[1], ".ex", ".in", 1)
	} else {
		filename = os.Args[1]
	}
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	monkeys := make([]Monkey, 0)

	for {
		m, ok := readMonkey(scanner)
		monkeys = append(monkeys, m)
		if !ok {
			break
		}
		fmt.Println("--------")
	}
	fmt.Println("--------")
	commonMultiplier := 1
	for j := range monkeys {
		commonMultiplier *= monkeys[j].testval
	}

	const ROUNDS = 10000
	for i := 0; i < ROUNDS; i++ {
		for j := range monkeys {
			debug("Monkey", j)
			m := &monkeys[j]
			for k := range m.items {
				debug("Item with worry level", m.items[k])
				if m.op == "+" {
					m.items[k] += m.factor
				} else if m.op == "*" {
					if m.factor == -1 {
						m.items[k] *= m.items[k]
					} else {
						m.items[k] *= m.factor
					}
				} else {
					panic("unknown op")
				}
				debug("-> after inspection", m.items[k])
				// monkey gets bored
				if B_SIDE {
					m.items[k] %= commonMultiplier
				} else {
					m.items[k] /= 3
				}
				debug("-> after bored", m.items[k])
				if m.items[k]%m.testval == 0 {
					debug("thrown to", m.trueTarget)
					monkeys[m.trueTarget].items = append(monkeys[m.trueTarget].items, m.items[k])
					debug("-> other after throw", monkeys[m.trueTarget].items)
				} else {
					debug("thrown to", m.falseTarget)
					monkeys[m.falseTarget].items = append(monkeys[m.falseTarget].items, m.items[k])
					debug("-> other after throw", monkeys[m.falseTarget].items)
				}
				m.activity++
			}
			m.items = m.items[:0]
		}
		debug("After round", i+1)
		for j := range monkeys {
			debug("Monkey", j, "items", monkeys[j].items)
		}
		debug("--------")
	}
	fmt.Println("--------")
	activites := make([]int, len(monkeys))
	for j := range monkeys {
		fmt.Println("Monkey", j, "activity", monkeys[j].activity)
		activites[j] = monkeys[j].activity
	}
	fmt.Println("--------")

	maxIndex, monkeyBusiness := max(activites)
	activites[maxIndex] = -1

	_, max := max(activites)
	monkeyBusiness *= max

	fmt.Println("Common multiplier", commonMultiplier)
	fmt.Println("--------")
	fmt.Println("Monkey business", monkeyBusiness)
}
