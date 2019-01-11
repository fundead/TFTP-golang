package tftp

type ConnectionService struct {
	FileStore     MemoryFileStore
	PendingReads  map[string][]byte
	PendingWrites map[string][]byte
}

// New instantiates a ConnectionService
func (cs ConnectionService) New() ConnectionService {
	cs.PendingReads = make(map[string][]byte)
	cs.PendingWrites = make(map[string][]byte)
	cs.FileStore = MemoryFileStore{}.New()
	return cs
}

func (cs ConnectionService) openRead(fileName string) []byte {
	allBytes, err := cs.FileStore.Read(fileName)
	if err != nil {
		cs.PendingReads[fileName] = allBytes
		return cs.readData(fileName)
	}
	return make([]byte, 0) // TODO
}

func (cs ConnectionService) readData(fileName string) []byte { // TODO by clientId
	// TODO bounds
	allData := cs.PendingReads[fileName]
	cs.PendingReads[fileName] = allData[516:]

	if len(allData) < 516 {
		delete(cs.PendingReads, fileName)
	}
	return allData[:516]
}

func (cs ConnectionService) openWrite(fileName string) {
	cs.PendingWrites[fileName] = make([]byte, 0)
}

func (cs ConnectionService) writeData(fileName string, data []byte) {
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
