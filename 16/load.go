package main

import (
	"encoding/json"
	"os"
)

func Load[T any](filename string) []T {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var arr []T
	err = json.Unmarshal(data, &arr)
	if err != nil {
		panic(err)
	}

	return arr
}
