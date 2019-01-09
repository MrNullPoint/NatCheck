package main

import (
	"NatCheck/stun"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println(stun.RunCheck(os.Args[1]))
	} else {
		fmt.Println(stun.RunCheck(""))
	}
}
