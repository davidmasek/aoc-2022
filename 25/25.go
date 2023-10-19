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

func Decimal(snafu string) int {
	sum := 0
	for i := 0; i < len(snafu); i++ {
		j := len(snafu) - i - 1
		c := string(snafu[j])
		var v int
		switch c {
		case "-":
			v = -1
		case "=":
			v = -2
		case "0":
			v = 0
		case "1":
			v = 1
		case "2":
			v = 2
		default:
			return -1
		}
		a := v * Pow(5, i)
		sum += a
	}
	return sum
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

	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		fmt.Println(line)
		sum := Decimal(line)

		if !test {
			fmt.Println(sum)
			fmt.Println("-----")
		}
		total += sum
	}
	fmt.Println("Decimal:", total)
	fmt.Println("SNAFU:", SNAFU(total))
	if !test {
		fmt.Println("-----")
		fmt.Println("SNAFU:", 1, SNAFU(1))
		fmt.Println("SNAFU:", 3, SNAFU(3))
		fmt.Println("SNAFU:", 4, SNAFU(4))
		fmt.Println("SNAFU:", 5, SNAFU(5))
		fmt.Println("SNAFU:", 10, SNAFU(10))
		fmt.Println("SNAFU:", 2022, SNAFU(2022))
		fmt.Println("SNAFU:", 12345, SNAFU(12345))
		fmt.Println("SNAFU:", 314159265, SNAFU(314159265))
	}

}

func SNAFU(x int) string {
	original := x
	// to base 5
	var digits []int = make([]int, 0)
	for x > 0 {
		d := x % 5
		x /= 5
		digits = append(digits, d)
	}
	// add a leading zero to avoid out of index error
	digits = append(digits, 0)
	// to SNAFU
	for i, v := range digits {
		if v == 3 {
			digits[i+1] += 1
			digits[i] = -2
		} else if v == 4 {
			digits[i+1] += 1
			digits[i] = -1
		} else if v == 5 {
			digits[i+1] += 1
			digits[i] = 0
		}
	}

	s := ""
	for _, v := range digits {
		var c string
		switch v {
		case -2:
			c = "="
		case -1:
			c = "-"

		default:
			c = fmt.Sprintf("%d", v)
		}
		s = c + s
	}
	if s[0] == '0' {
		s = s[1:]
	}
	if Decimal(s) != original {
		fmt.Println("    ERROR---v", s, Decimal(s))
	}
	return s
}

func Pow(x, n int) int {
	if n < 0 {
		panic("not implemented for native powers")
	}
	if n == 0 {
		return 1
	}
	y := 1
	for n > 1 {
		// n is odd
		if n&1 == 1 {
			y = x * y
			n--
		}
		x = x * x
		n = n / 2
	}
	return x * y
}
