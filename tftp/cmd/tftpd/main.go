package main

import (
	"fmt"
	"log"

	tftp "igneous.io/tftp/lib"
)

func main() {
	cs := tftp.ConnectionService{}.New()
	go tftp.Listen(&cs)
	log.Println("Press a key to exit")
	fmt.Scanln()
}
