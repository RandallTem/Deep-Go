package main

import (
	"fmt"
	"iter"
)

func Range() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < 10; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func main() {
	for value := range Range() {
		if value == 5 {
			return
		}
		fmt.Println(value)
	}
}
