package tftp

import "errors"

type MemoryFileStore struct {
	files map[string]File
}

type File struct {
	fileName string
	data     []byte
}

func (mem MemoryFileStore) Read(fileName string) ([]byte, error) {
	if file, ok := mem.files[fileName]; ok {
		return file.data, nil
	}
	return nil, errors.New("Tried to read filename which doesn't exist in mem")
}

func (mem MemoryFileStore) Write(fileName string, data []byte) bool {
	if currentFile, ok := mem.files[fileName]; ok {
		currentBytes := currentFile.data
		currentFile.data = append(currentBytes, data...)
		mem.files[fileName] = currentFile
	} else {
		file := File{fileName, data}
		mem.files[fileName] = file
	}
	return true
}
