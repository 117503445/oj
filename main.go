package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Println("1")
	s := flag.String("s", "", "source code path")
	fmt.Println(*s)
	fmt.Println("1")

}
