package main

import (
	"fmt"
	"iter"
)

type Node struct {
	Value int
	Next  *Node
}

func IterateList(head *Node) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		counter := 0
		for curr := head; curr != nil; curr = curr.Next {
			counter++
			if !yield(curr.Value, counter) {
				return
			}
		}
	}
}

func main() {
	list := &Node{Value: 1}
	list.Next = &Node{Value: 2}
	list.Next.Next = &Node{Value: 3}
	list.Next.Next.Next = &Node{Value: 4}

	for v, c := range IterateList(list) {
		fmt.Println(c, ":", v) // 1 : 1, 2 : 2, 3 : 3
		if c == 3 {
			break
		}
	}
}
