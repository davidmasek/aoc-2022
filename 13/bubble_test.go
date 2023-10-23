package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	arr := []int{4, 3, 1, 2, 2}
	t.Log("input:", arr)
	BubbleSort(arr, func(a, b *int) bool { return *a < *b })
	for i := range arr {
		if i+1 == len(arr) {
			break
		}
		if arr[i+1] < arr[i] {
			t.Error(fmt.Sprintf("[pos %d]", i), arr[i], ">", arr[i+1])
		}
	}
	t.Log("sorted:", arr)

	arr = make([]int, 1000)
	for i := range arr {
		arr[i] = rand.Intn(100)
	}
	BubbleSort(arr, func(a, b *int) bool { return *a < *b })
	for i := range arr {
		if i+1 == len(arr) {
			break
		}
		if arr[i+1] < arr[i] {
			t.Error(fmt.Sprintf("[pos %d]", i), arr[i], ">", arr[i+1])
		}
	}

	arrFloat := make([]float32, 1000)
	for i := range arrFloat {
		arrFloat[i] = rand.Float32()
	}
	BubbleSort(arrFloat, func(a, b *float32) bool { return *a < *b })
	for i := range arrFloat {
		if i+1 == len(arrFloat) {
			break
		}
		if arrFloat[i+1] < arrFloat[i] {
			t.Error(fmt.Sprintf("[pos %d]", i), arrFloat[i], ">", arrFloat[i+1])
		}
	}
}
