package main

// Return true if a < b, false if a >= b.
type less[V any] func(a, b *V) bool

func BubbleSort[V any](arr []V, less less[V]) []V {
	for {
		sorted := true
		for i := range arr {
			if i+1 == len(arr) {
				break
			}
			if less(&arr[i+1], &arr[i]) {
				arr[i], arr[i+1] = arr[i+1], arr[i]
				sorted = false
			}
		}
		if sorted {
			break
		}
	}
	return arr
}
