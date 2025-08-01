package main

import (
	"fmt"
	"maps"
	"math"
)

func main() {
	m := map[string]float64{
		"pi": math.Pi,
		"e":  math.E,
	}

	for k, v := range maps.All(m) {
		fmt.Println("const", k, "val", v) // const pi val 3.14..., const e val 2.71...
	}

	for k := range maps.Keys(m) {
		fmt.Println("key", k) // key e, key pi
	}

	for v := range maps.Values(m) {
		fmt.Println("value", v) // value 3.14..., value 2.71
	}
}
