package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Print("err args num >1 ...")
		os.Exit(1)
	}

	fmt.Print("ok")
}
