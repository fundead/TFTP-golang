package main

import (
	"fmt"
	"log"

	tftp "igneous.io/tftp/lib"
)

func main() {
	go tftp.Listen(&tftp.ConnectionService{})
	log.Println("Press a key to exit")
	fmt.Scanln()
}
