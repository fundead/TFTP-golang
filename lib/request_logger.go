package tftp

import "os"

const logFileName = "requestLog.txt"

func LogRequest(text string) {
	f, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text + "\n"); err != nil {
		panic(err)
	}
}

func LogWriteRequest(fileName string) {
	LogRequest("Got write request for: " + fileName)
}

func LogReadRequest(fileName string) {
	LogRequest("Got read request for: " + fileName)
}

func LogFileNotFound(fileName string) {
	LogRequest("Got a read request for file which isn't in memory: " + fileName)
}

func LogUnknownRequest() {
	LogRequest("Got unknown request (unknown opcode)")
}
