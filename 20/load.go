package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ReadInts(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var nums []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" { // skip any empty lines
			num, err := strconv.Atoi(line)
			if err != nil {
				return nil, err
			}
			nums = append(nums, num)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nums, nil
}
