package main

import (
	"fmt"

	"github.com/nykolynoleg/goptional"
)

func main() {
	// Create a pair of Optionals.
	pair := goptional.Pair[goptional.Optional[int], goptional.Optional[string]]{
		First:  goptional.Of(123),
		Second: goptional.Of("gm"),
	}

	// Unwrap the given optional pair.
	// Return two empty optionals if the given optional is empty.
	opt1, opt2 := goptional.Unzip(goptional.Of(&pair))

	fmt.Println(opt1.Get()) // 123
	fmt.Println(opt2.Get()) // gm

	// Create empty pair.
	emptyPair := goptional.Empty[*goptional.Pair[goptional.Optional[int], goptional.Optional[string]]]()
	opt1, opt2 = goptional.Unzip(emptyPair)

	fmt.Println(opt1.IsEmpty()) // true
	fmt.Println(opt2.IsEmpty()) // true
}
