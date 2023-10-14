package main

import (
	"fmt"
	"os"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

func check(err error) {
	if err != nil {
		panic(err)
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
	i := 0
	var counter map[rune]int
	for {
		if i == 0 {
			counter = make(map[rune]int)
		}
		var line string
		n, _ := fmt.Fscanln(f, &line)
		if n != 1 {
			break
		}
		used := mapset.NewSet[rune]()
		for _, a := range line {
			used.Add(a)
		}
		for a := range used.Iter() {
			counter[a]++
		}
		i++
		if i == 3 {
			i = 0
			// TODO
			// fmt.Println(counter)
			for k, v := range counter {
				if v == 3 {
					var keyValue int
					if k >= 'a' && k <= 'z' {
						keyValue = int(k - 'a' + 1)
					} else {
						keyValue = int(k - 'A' + 27)
					}
					fmt.Println(k, keyValue)
					sum += keyValue
				}
			}
		}
	}
	fmt.Println(sum)
}
