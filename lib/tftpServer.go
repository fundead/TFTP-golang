package tftp

import (
	"log"
	"net"
)

type UDPPacket struct {
	Address net.Addr
	Data    []byte
}

// Listen establishes a new TFTPServer and listens on port 69
func Listen(connectionService *ConnectionService) {
	inbound := make(chan UDPPacket)
	packetConn, err := net.ListenPacket("udp", "127.0.0.1:69")
	if err != nil {
		log.Fatal("Error listening:", err)
		return
	}
	defer packetConn.Close()

	go read(packetConn, inbound)
	go process(packetConn, inbound, connectionService)
	log.Println("Server listening on port 69")
	for {
	}
}

func read(conn net.PacketConn, inbound chan UDPPacket) {
	for {
		buffer := make([]byte, MaxPacketSize)
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Println("Error: UDP packet read", err)
			continue
		}
		udpPacket := UDPPacket{addr, buffer[:n]}
		inbound <- udpPacket
	}
}

func process(pc net.PacketConn, inbound chan UDPPacket, connectionService *ConnectionService) {
	for {
		select {
		case udpPacket := <-inbound:
			packet, err := ParsePacket(udpPacket.Data)
			if err != nil {
				log.Println(err)
			}
			parse(pc, udpPacket.Address, packet, connectionService)
		}
	}
}

func parse(pc net.PacketConn, addr net.Addr, packet Packet, connectionSvc *ConnectionService) {
	switch packet.(type) {
	case *PacketRequest:
		handleRequest(pc, addr, packet.(*PacketRequest), connectionSvc)
	case *PacketAck:
		handleAck(pc, addr, packet.(*PacketAck), connectionSvc)
	case *PacketData:
		handleData(pc, addr, packet.(*PacketData), connectionSvc)
	case *PacketError:
		log.Println("Warning: received an error packet")
	default:
		LogUnknownRequest()
	}
}

// For new read/write requests: send appropriate ACK or DATA response
func handleRequest(pc net.PacketConn, addr net.Addr, pr *PacketRequest, connectionSvc *ConnectionService) {
	if pr.Op == OpRRQ { // Read Request
		LogReadRequest(pr.Filename)
		data, err := connectionSvc.openRead(addr.String(), pr.Filename)
		if err != nil {
			LogFileNotFound(pr.Filename)
			sendResponse(pc, addr, &PacketError{0x1, "File not found (error opening file read)"})
		} else {
			sendResponse(pc, addr, &PacketData{0x1, data})
		}
	} else if pr.Op == OpWRQ { // Write Request
		LogWriteRequest(pr.Filename)
		connectionSvc.openWrite(addr.String(), pr.Filename)
		sendResponse(pc, addr, &PacketAck{0})
	}
}

// For a read: sends the next DATA block in response to an ACK
func handleAck(pc net.PacketConn, addr net.Addr, pa *PacketAck, connectionSvc *ConnectionService) {
	payload := connectionSvc.readData(addr.String(), pa.BlockNum+1)
	dataPacket := &PacketData{pa.BlockNum + 1, payload}
	sendResponse(pc, addr, dataPacket)
}

// For a write: sends an ACK in response to a DATA payload
func handleData(pc net.PacketConn, addr net.Addr, pd *PacketData, connectionSvc *ConnectionService) {
	connectionSvc.writeData(addr.String(), pd.BlockNum, pd.Data)
	ackPacket := &PacketAck{pd.BlockNum}
	sendResponse(pc, addr, ackPacket)
}

func sendResponse(pc net.PacketConn, addr net.Addr, p Packet) {
	_, err := pc.WriteTo(p.Serialize(), addr)
	if err != nil {
		log.Fatalln("Error: failed to send response to packet", err)
	}
}
