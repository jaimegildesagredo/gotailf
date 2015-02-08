package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jaimegildesagredo/tailf"
)

func main() {
	path, numLines := parseArgs()
	output := make(chan string, 1)
	var line string

	go tailf.Tailf(path, numLines, output)

	for {
		select {
		case line = <-output:
			fmt.Println(line)
		}
	}
}

func parseArgs() (string, int) {
	numLines := flag.Int("n", 1, "Output the last N lines")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("Usage: tailf [options] file.txt")

	}

	return flag.Args()[0], *numLines
}
