package main

import (
	"fmt"
	"slices"
)

func main() {
	s := []int{1, 2}

	for i, v := range slices.All(s) { // returns iter.Seq2
		fmt.Println("item", i, "is", v) // item 0 is 1, item 1 is 2
	}

	for v := range slices.Values(s) { // returns iter.Seq
		fmt.Println(v) // 1 2
	}
}
