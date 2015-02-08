package tailf

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"time"
)

const (
	NEWLINE = "\n"
)

func Tailf(path string, numLines int, output chan string) {
	var file *os.File
	var fileContent []byte
	var fileStat os.FileInfo
	var readOffset int64 = getOffsetBytesForLastLines(path, numLines)
	var err error
	var linesScanner *bufio.Scanner

	newLineBytes := []byte(NEWLINE)
	newLineBytesLen := int64(len(newLineBytes))

	for {
		file, err = os.Open(path)

		if err != nil {
			log.Println("Error opening file", err)
			continue
		}

		fileStat, err = file.Stat()

		if err != nil {
			log.Println("Error getting file stat", err)
			continue
		}

		fileContent = make([]byte, fileStat.Size()-readOffset)

		_, err = file.ReadAt(fileContent, readOffset)

		if err != nil {
			log.Println("Error reading file content", err)
			continue
		}

		file.Close()

		if len(fileContent) > 0 {
			linesScanner = bufio.NewScanner(bytes.NewReader(fileContent))

			for linesScanner.Scan() {
				readOffset += int64(len(linesScanner.Bytes())) + newLineBytesLen
				output <- string(linesScanner.Text())
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func getOffsetBytesForLastLines(path string, numLines int) int64 {
	var file *os.File
	var chunk []byte
	var fileStat os.FileInfo
	var readOffset int64
	var err error

	file, err = os.Open(path)
	defer file.Close()

	if err != nil {
		log.Println("Error opening file", err)
		return 0
	}

	fileStat, err = file.Stat()

	if err != nil {
		log.Println("Error getting file stat", err)
		return 0
	}

	newLineBytes := []byte(NEWLINE)
	newLineBytesLen := int64(len(newLineBytes))
	iteration := 0
	linesRead := 0

	for {
		chunk = make([]byte, newLineBytesLen)

		iteration += 1
		readOffset = fileStat.Size() - newLineBytesLen*int64(iteration)
		_, err = file.ReadAt(chunk, readOffset)

		if err != nil {
			log.Println("Error reading file content", err)
			continue
		}

		if bytes.Compare(chunk, newLineBytes) == 0 {
			linesRead += 1
		}

		if linesRead >= numLines {
			return readOffset
		}
	}
}
