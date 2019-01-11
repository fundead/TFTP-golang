package tftp

import (
	"errors"
)

type MemoryFileStore struct {
	Files map[string]File
}

type File struct {
	FileName string
	Data     []byte
}

func (mem MemoryFileStore) New() MemoryFileStore {
	mem.Files = make(map[string]File)
	return mem
}

func (mem MemoryFileStore) Read(fileName string) ([]byte, error) {
	if file, ok := mem.Files[fileName]; ok {
		return file.Data, nil
	}
	return nil, errors.New("Tried to read filename which doesn't exist in mem")
}

func (mem MemoryFileStore) Write(fileName string, data []byte) bool {
	if currentFile, ok := mem.Files[fileName]; ok {
		currentBytes := currentFile.Data
		currentFile.Data = append(currentBytes, data...)
		mem.Files[fileName] = currentFile
	} else {
		file := File{fileName, data}
		mem.Files[fileName] = file
	}
	return true
}
