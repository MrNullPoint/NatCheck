package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println(RunCheck(os.Args[1]))
	} else {
		fmt.Println(RunCheck(""))
	}
}
