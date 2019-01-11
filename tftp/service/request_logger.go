package tftp

import "os"

const logFileName = "log.txt"

func LogRequest(text string) {
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func LogWriteRequest(fileName string) {
	LogRequest("write for : " + fileName)
}

func LogReadRequest(fileName string) {
	LogRequest("read for: " + fileName)
}
