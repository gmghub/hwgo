package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	s := "Hello, OTUS!"
	res := stringutil.Reverse(s)
	fmt.Println(res)
}
