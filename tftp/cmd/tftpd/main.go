package main

import (
	"fmt"
	"log"

	tftp "igneous.io/tftp/lib"
)

func main() {
	cs := tftp.ConnectionService{}.New()
	go tftp.Listen(&cs)
	log.Println("Press enter to exit")
	fmt.Scanln()
}
