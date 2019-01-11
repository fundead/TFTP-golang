package tftp

import "log"

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

func (cs ConnectionService) openRead(address string, fileName string) []byte {
	allBytes, err := cs.FileStore.Read(fileName)
	if err != nil {
		log.Fatalln("Error on openRead for file "+fileName, err)
	}
	cs.PendingReads[fileName] = allBytes
	return cs.readData(fileName)
}

func (cs ConnectionService) readData(address string, blockNum uint16) []byte {
	remainingData := cs.PendingReads[fileName]
	if len(remainingData) < 516 {
		delete(cs.PendingReads, fileName)
		return remainingData
	}
	cs.PendingReads[fileName] = remainingData[516:]
	return remainingData[:516]
}

func (cs ConnectionService) openWrite(address string, fileName string) {
	cs.PendingWrites[fileName] = make([]byte, 0)
}

func (cs ConnectionService) writeData(address string, blockNum uint16, data []byte) {
	allData := cs.PendingWrites[fileName]
	allData = append(allData, data...)
	cs.PendingWrites[fileName] = allData
	if len(data) < 516 {
		cs.closeWrite(fileName)
	}
}

func (cs ConnectionService) closeWrite(fileName string) {
	cs.FileStore.Write(fileName, cs.PendingWrites[fileName])
	delete(cs.PendingWrites, fileName)
}
