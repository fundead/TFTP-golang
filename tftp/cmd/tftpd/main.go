package main

import (
	"fmt"

	tftp "igneous.io/tftp/lib"
)

func main() {
	go tftp.Listen(&tftp.ConnectionService{
		FileStore: tftp.MemoryFileStore{},
	})
	fmt.Scanln() // TODO exit ctrl-c
}
