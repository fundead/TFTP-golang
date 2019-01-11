package tftp

import (
	"log"
	"net"
	"reflect"
)

type UDPPacket struct {
	Address *net.UDPAddr
	Data    []byte
}

func Listen(connectionService *ConnectionService) {
	inbound := make(chan UDPPacket) // TODO add buffered channel length
	addr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 6900,
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal("Error listening:", err)
		return
	}
	defer conn.Close()

	go read(conn, inbound)
	go process(inbound, connectionService)
	log.Println("Server listening on port 6900")

	// TODO
	for {
	}
}

func read(conn *net.UDPConn, inbound chan UDPPacket) {
	for {
		buffer := make([]byte, MaxPacketSize)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error: UDP packet read", err)
			continue
		}
		udpPacket := UDPPacket{addr, buffer[:n]}
		inbound <- udpPacket
	}
}

func process(inbound chan UDPPacket, connectionService *ConnectionService) {
	for {
		select {
		case udpPacket := <-inbound:
			packet, err := ParsePacket(udpPacket.Data)
			if err != nil {
				log.Println(err)
			}
			log.Println(reflect.TypeOf(packet)) // TODO rm
			parse(udpPacket.Address, packet, connectionService)
		}
	}
}

func parse(addr *net.UDPAddr, packet Packet, connectionSvc *ConnectionService) {
	switch packet.(type) {
	case *PacketRequest:
		handleRequest(addr, packet.(*PacketRequest), connectionSvc)
	case *PacketAck:
		handleAck(addr, packet.(*PacketAck), connectionSvc)
	case *PacketData:
		handleData(addr, packet.(*PacketData), connectionSvc)
	case *PacketError:
		log.Println("Packet Error")
	default:
		log.Println("Unknown packet/opcode received")
	}
}

func handleRequest(addr *net.UDPAddr, pr *PacketRequest, connectionSvc *ConnectionService) {
	if pr.Op == OpRRQ { // Read Request
		data := connectionSvc.openRead(pr.Filename)
		go sendResponse(addr, &PacketData{0x1, data})
	} else if pr.Op == OpWRQ { // Write Request
		connectionSvc.openWrite(pr.Filename)
		go sendResponse(addr, &PacketAck{0})
	}
	log.Println("Opened request with type " + string(pr.Op))
}

func sendResponse(addr *net.UDPAddr, p Packet) {
	serverAddr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 6900,
	}
	conn, err := net.DialUDP("udp", &serverAddr, addr)
	if err != nil {
		conn.WriteToUDP(p.Serialize(), addr)
	} else {
		log.Fatalln("Error: failed to write next data in response to ACK")
	}
}

// For a read: sends the next DATA block in response to an ACK
func handleAck(addr *net.UDPAddr, pa *PacketAck, connectionSvc *ConnectionService) {
	payload := connectionSvc.readData("") //pa.BlockNum)
	dataPacket := &PacketData{pa.BlockNum + 1, payload}
	sendResponse(addr, dataPacket)
}

// For a write: sends an ACK in response to a DATA payload
func handleData(addr *net.UDPAddr, pd *PacketData, connectionSvc *ConnectionService) {
	//connectionSvc.writeData(pd.BlockNum, pd.Data)
	connectionSvc.writeData("", pd.Data) // TODO
	ackPacket := &PacketAck{pd.BlockNum}
	sendResponse(addr, ackPacket)
}
