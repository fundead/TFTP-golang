package tftp

type MemoryFileStore struct {
	files map[string]File
}

type File struct {
	fileName   string
	data       []byte
	isComplete bool
}

func (mem MemoryFileStore) Read(file string) []byte {
	// TODO guard
	return mem.files[file].data
}

func (mem MemoryFileStore) Write(f File) bool {
	if val, ok := mem.files[f.fileName]; ok {
		currentBytes := mem.files[f.fileName].data
		f.data = append(currentBytes, f.data...)
		mem.files[f.fileName] = f
	} else {
		mem.files[f.fileName] = f
	}
}
