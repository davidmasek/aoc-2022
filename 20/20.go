package main

import (
	"fmt"
)

const KEY = 811589153

type Leaf struct {
	Value            int
	Position         int
	OriginalPosition int
	Next             *Leaf
	Prev             *Leaf
}

func (leaf *Leaf) MoveLeft(modulo int) {
	leaf.Position = (leaf.Position - 1 + modulo) % modulo
	leaf.Prev.Position = (leaf.Prev.Position + 1 + modulo) % modulo

	leaf.Next.Prev = leaf.Prev
	leaf.Prev.Next = leaf.Next
	leaf.Next = leaf.Prev
	leaf.Prev = leaf.Prev.Prev
	leaf.Next.Prev = leaf
	leaf.Prev.Next = leaf
}

func (leaf *Leaf) MoveRight(modulo int) {
	leaf.Position = (leaf.Position + 1 + modulo) % modulo
	leaf.Next.Position = (leaf.Next.Position - 1 + modulo) % modulo

	leaf.Prev.Next = leaf.Next
	leaf.Next.Prev = leaf.Prev
	leaf.Prev = leaf.Next
	leaf.Next = leaf.Next.Next
	leaf.Prev.Next = leaf
	leaf.Next.Prev = leaf
}

func (leaf *Leaf) FindHead() *Leaf {
	head := leaf
	for head.Position != 0 {
		head = head.Prev
	}
	return head
}

func (leaf *Leaf) Print() {
	head := leaf.FindHead()
	current := head
	for {
		fmt.Printf("%d, ", current.Value)
		current = current.Next
		if current == head {
			break
		}
	}
	fmt.Println()
	fmt.Println()
}

func solveLinked(arr []int) {
	originalOrder := make([]*Leaf, 0, len(arr))
	head := &Leaf{arr[0], 0, 0, nil, nil}
	originalOrder = append(originalOrder, head)
	prev := head
	for i := 1; i < len(arr); i++ {
		leaf := &Leaf{
			Value:            arr[i],
			Position:         i,
			OriginalPosition: i,
			Next:             nil,
			Prev:             prev,
		}
		originalOrder = append(originalOrder, leaf)
		prev.Next = leaf
		prev = leaf
	}
	last := prev
	last.Next = head
	head.Prev = last

	MOD := len(arr) - 1
	N_ITERS := 10
	for i := 0; i < N_ITERS; i++ {
		for _, leaf := range originalOrder {
			shift := leaf.Value * (KEY % MOD)
			for shift > 0 {
				shift %= MOD
				leaf.MoveRight(len(arr))
				shift--
			}
			for shift < 0 {
				shift %= MOD
				leaf.MoveLeft(len(arr))
				shift++
			}
			// fmt.Println(leaf.Value, "moves between", leaf.Prev.Value, "and", leaf.Next.Value)
			// leaf.Print()
		}
	}

	leaf := originalOrder[0]
	for leaf.Value != 0 {
		leaf = leaf.Next
	}

	sum := 0
	for i := 0; i < 3000; i++ {
		leaf = leaf.Next
		j := i + 1
		if j%1000 == 0 {
			fmt.Println(j, leaf.Value, leaf.Value*KEY)
			sum += leaf.Value * KEY
		}
	}
	fmt.Println(sum)
}

func main() {
	arr, err := ReadInts("ex.txt")
	if err != nil {
		panic(err)
	}
	solveLinked(arr)
	arr, err = ReadInts("in.txt")
	if err != nil {
		panic(err)
	}
	solveLinked(arr)
}
