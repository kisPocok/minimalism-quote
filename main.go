package main

import (
	"fmt"

	"minimal/grabber/clients"

	"golang.org/x/sys/unix"
)

func main() {
	quote, err := clients.Minimalmaxism().GrabQuote()
	if err != nil {
		unix.Exit(1)
	}
	fmt.Println(quote)
	unix.Exit(0)
}
