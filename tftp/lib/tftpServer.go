package tftp

import (
	"fmt"
	"log"
	"net"
	"reflect"
)

type UDPPacket struct {
	Address net.Addr
	Data    []byte
}

func Listen(fileStore *MemoryFileStore) {
	inbound := make(chan UDPPacket) // TODO add buffered channel length
	conn, err := net.ListenPacket("udp", ":6900")
	if err != nil {
		log.Fatal("Error listening:", err)
		return
	}
	defer conn.Close()

	go read(conn, inbound)
	go process(inbound)
	fmt.Println("Server listening on port 6900")

	for {
	}
}

func read(conn net.PacketConn, inbound chan UDPPacket) {
	for {
		buffer := make([]byte, MaxPacketSize)
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Println("ERROR: UDP read", err)
			continue
		}
		udpPacket := UDPPacket{addr, buffer[:n]}
		inbound <- udpPacket
	}
}

func process(inbound chan UDPPacket) {
	for {
		select {
		case udpPacket := <-inbound:
			packet, err := ParsePacket(udpPacket.Data)
			if err != nil {
				log.Println(err)
			}
			log.Println(reflect.TypeOf(packet))

			parse(packet)
		}
	}
}

func parse(packet Packet) {
	switch packet.(type) {
	case *PacketRequest:
		log.Println("RRQ/WRQ")
	case *PacketAck:
		log.Println("Packet ACK")
	case *PacketData:
		log.Println("Data")
	case *PacketError:
		log.Println("Packet Error")
	default:
		log.Println("Unknown packet/opcode received")
	}
}
