package tftp

import (
	"log"
)

type Connection struct {
	Address  string
	BlockNum uint16
}

type ConnectionService struct {
	FileStore     MemoryFileStore
	PendingReads  map[Connection]File
	PendingWrites map[Connection]File
}

// New instantiates a ConnectionService
func (cs ConnectionService) New() ConnectionService {
	cs.PendingReads = make(map[Connection]File)
	cs.PendingWrites = make(map[Connection]File)
	cs.FileStore = MemoryFileStore{}.New()
	return cs
}

func (cs ConnectionService) openRead(address string, fileName string) ([]byte, error) {
	connection := Connection{address, 1}
	allBytes, err := cs.FileStore.Read(fileName)
	if err != nil {
		log.Println("Got request for file which doesn't exist in memory: " + fileName)
		return nil, err
	}
	file := File{fileName, allBytes}
	cs.PendingReads[connection] = file
	return cs.readData(address, 1), nil
}

func (cs ConnectionService) readData(address string, blockNum uint16) []byte {
	connection := Connection{address, blockNum}
	file := cs.PendingReads[connection]
	remainingData := file.Data
	if len(remainingData) <= 512 {
		delete(cs.PendingReads, connection)
		return remainingData
	}
	file.Data = remainingData[512:]
	delete(cs.PendingReads, connection)
	connection.BlockNum++
	cs.PendingReads[connection] = file
	return remainingData[:512]
}

func (cs ConnectionService) openWrite(address string, fileName string) {
	connection := Connection{address, 1}
	cs.PendingWrites[connection] = File{fileName, make([]byte, 0)}
}

func (cs ConnectionService) writeData(address string, blockNum uint16, data []byte) {
	connection := Connection{address, blockNum}
	file := cs.PendingWrites[connection]
	currentData := file.Data
	file.Data = append(currentData, data...)
	delete(cs.PendingWrites, connection)
	connection.BlockNum++
	cs.PendingWrites[connection] = file
	if len(data) < 512 {
		cs.closeWrite(connection)
	}
}

func (cs ConnectionService) closeWrite(connection Connection) {
	file := cs.PendingWrites[connection]
	cs.FileStore.Write(file.FileName, file.Data)
	delete(cs.PendingWrites, connection)
}
