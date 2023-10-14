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
		var line string
		n, _ := fmt.Fscanln(f, &line)
		if n != 1 {
			break
		}
		l, r := line[:len(line)/2], line[len(line)/2:]
		var v int
		for _, a := range l {
			for _, b := range r {
				if a == b {
					if a >= 'a' && a <= 'z' {
						v = int(a - 'a' + 1)
					} else {
						v = int(a - 'A' + 27)
					}
				}
			}
		}
		fmt.Println(l, r, v)
		sum += v
	}
	fmt.Println(sum)
}
