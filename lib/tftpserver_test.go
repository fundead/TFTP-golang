package tftp

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func SetupTestServer(port int) (cs ConnectionService, ch chan UDPPacket, packetConn net.PacketConn) {
	packetConn, _ = net.ListenPacket("udp", "127.0.0.1:"+strconv.Itoa(port))
	cs = ConnectionService{}.New()
	ch = make(chan UDPPacket)
	return cs, ch, packetConn
}

func TestReadFile(t *testing.T) {
	// Arrange
	connSvc, ch, packetConn := SetupTestServer(5001)
	go process(packetConn, ch, &connSvc)
	connSvc.FileStore.Files["testreadfile"] = File{"testreadfile", make([]byte, 513)}

	// Act
	sourceAddress := &net.UDPAddr{net.ParseIP("127.0.0.1"), 10001, ""}
	readRequest := &PacketRequest{0x1, "testreadfile", "octet"}
	packet := UDPPacket{sourceAddress, readRequest.Serialize()}
	ch <- packet

	// Assert
	time.Sleep(20 * time.Millisecond)
	if len(connSvc.PendingReads) != 1 {
		t.Fatal("Expected read not in flight")
	}
}

func TestWriteFile(t *testing.T) {
	// Arrange
	connSvc, ch, packetConn := SetupTestServer(5002)
	go read(packetConn, ch)
	go process(packetConn, ch, &connSvc)

	// Act
	sourceAddress := &net.UDPAddr{net.ParseIP("127.0.0.1"), 5002, ""}
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

// A smoke test for confirming the instance/channel communication works without error
func TestInstatiation(t *testing.T) {
	connSvc, ch, packetConn := SetupTestServer(5000)
	go process(packetConn, ch, &connSvc)
	ack := &PacketAck{0x1}
	sourceAddress := &net.UDPAddr{net.ParseIP("127.0.0.1"), 10000, ""}
	packet := UDPPacket{sourceAddress, ack.Serialize()}
	ch <- packet
}
