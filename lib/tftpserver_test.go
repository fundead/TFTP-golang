package tftp

import (
	"net"
	"testing"
	"time"
)

func SetupTestServer() (cs ConnectionService, ch chan UDPPacket, packetConn net.PacketConn) {
	packetConn, _ = net.ListenPacket("udp", "127.0.0.1:10000")
	cs = ConnectionService{}.New()
	ch = make(chan UDPPacket)
	return cs, ch, packetConn
}

func TestWriteFile(t *testing.T) {
	// Arrange
	connSvc, ch, packetConn := SetupTestServer()
	go read(packetConn, ch)
	go process(packetConn, ch, &connSvc)

	// Act
	sourceAddress := &net.UDPAddr{net.ParseIP("127.0.0.1"), 10000, ""}
	writeRequest := &PacketRequest{0x2, "testfile", "octet"}
	dataPacket := &PacketData{0x1, []byte("content")}
	sendResponse(packetConn, sourceAddress, writeRequest)
	sendResponse(packetConn, sourceAddress, dataPacket)

	// Assert
	time.Sleep(20 * time.Millisecond)
	file, err := connSvc.FileStore.Read("testfile")
	if file == nil || err != nil {
		t.Fatal("File wasn't written to MemoryFileStore", err)
	}
}

func TestReadFile(t *testing.T) {
	// Arrange
	connSvc, ch, packetConn := SetupTestServer()
	go process(packetConn, ch, &connSvc)
	connSvc.FileStore.Files["testreadfile"] = File{"testreadfile", make([]byte, 513)}

	// Act
	sourceAddress := &net.UDPAddr{net.ParseIP("127.0.0.1"), 10000, ""}
	readRequest := &PacketRequest{0x1, "testreadfile", "octet"}
	packet := UDPPacket{sourceAddress, readRequest.Serialize()}
	ch <- packet

	// Assert
	time.Sleep(20 * time.Millisecond)
	if len(connSvc.PendingReads) != 1 {
		t.Fatal("Expected read not in flight")
	}
}

func TestInstatiation(t *testing.T) {
	connSvc, ch, packetConn := SetupTestServer()
	go process(packetConn, ch, &connSvc)
	ack := &PacketAck{0x1}
	sourceAddress := &net.UDPAddr{net.ParseIP("127.0.0.1"), 10000, ""}
	packet := UDPPacket{sourceAddress, ack.Serialize()}
	ch <- packet
}
