package main

import (
	"fmt"

	tftp "igneous.io/tftp/lib"
)

func main() {
	go tftp.Listen()
	fmt.Scanln() // TODO exit ctrl-c
}
