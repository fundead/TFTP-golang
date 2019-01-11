package tftp

import "testing"

func MakeInbound() chan UDPPacket {
	return make(chan UDPPacket)
}

func BuildPacket(opcode byte) UDPPacket {
	b := []byte{0, opcode}
	u := UDPPacket{nil, b}
	return u
}

func BuildRequestPacket(opcode uint16, fileName string) UDPPacket {
	pr := PacketRequest{opcode, fileName, "octet"}
	return UDPPacket{nil, pr.Serialize()}
}

func BuildDataPacket(blockNumber uint16, data []byte) UDPPacket {
	pd := PacketData{blockNumber, []byte(data)}
	return UDPPacket{nil, pd.Serialize()}
}

func BuildAckPacket(blockNumber uint16) UDPPacket {
	pa := PacketAck{blockNumber}
	return UDPPacket{nil, pa.Serialize()}
}

func TestProcess(t *testing.T) {
	connSvc := ConnectionService{}
	ch := make(chan UDPPacket)
	go process(ch, &connSvc)
	packet := BuildPacket(4)
	ch <- packet
}

func TestWriteFile(t *testing.T) {
	writeRequest := BuildRequestPacket(OpWRQ, "testfile")
	//ack := BuildAckPacket(0x0) // assert receive this
	dataPacket := BuildDataPacket(0x1, []byte("content"))
	//ackData := BuildAckPacket(0x1) // assert receive this

	connSvc := ConnectionService{}
	ch := make(chan UDPPacket)
	go process(ch, &connSvc)

	ch <- writeRequest
	//ch <- ack
	ch <- dataPacket
	//ch <- ackData

	// assert file exists
}

func TestReadFile(t *testing.T) {
	readRequest := BuildRequestPacket(OpRRQ, "testfile")
	ackData := BuildAckPacket(0x1)

	connSvc := ConnectionService{}
	ch := make(chan UDPPacket)
	go process(ch, &connSvc)

	ch <- readRequest
	// assert receive data
	//dataPacket := BuildDataPacket(0x1, []byte("content")) // assert receive this
	ch <- ackData

	if 1 != 0 {
		t.Fatal("Not expected")
	}
}
