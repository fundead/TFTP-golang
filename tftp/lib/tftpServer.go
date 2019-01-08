package tftp

import (
	"fmt"
	"log"
	"net"
)

func Listen() {
	//inbound := make(chan Packet) // buffer_ch
	inbound := make(chan []byte) // buffer_ch
	conn, err := net.ListenPacket("udp", ":6969")
	if err != nil {
		log.Fatal("Error listening:", err)
		return
	}
	defer conn.Close()

	go read(conn, inbound)
	go process(inbound)
	fmt.Println("Server listening on port 69")

	for {
	}
	log.Println("Done")
}

func read(conn net.PacketConn, inbound chan []byte) {
	for {
		buffer := make([]byte, MaxPacketSize)
		//n, addr, err := conn.ReadFromUDP(buffer)
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Println("ERROR: UDP read", err)
			continue
		}
		//inbound <- Packet{addr, buffer[:n]}
		inbound <- buffer[:n]
	}
}

func process(inbound chan []byte) {
	for {
		select {
		case b := <-inbound:
			log.Println(b)
		}
	}
}
