package main

import (
	"fmt"
)

func main() {
	i := 0

	for i < 5 {
		fmt.Printf("%v. Hello World\n", i)
		i ++
	}
}