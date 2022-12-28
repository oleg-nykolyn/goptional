package main

import (
	"fmt"

	"github.com/nykolynoleg/goptional"
)

func main() {
	opt := goptional.Empty[int]()
	numAsJSON := "123"

	// Populate opt with the given JSON.
	err := opt.UnmarshalJSON([]byte(numAsJSON))

	fmt.Println(err == nil) // true
	fmt.Println(opt.Get())  // 123
}
