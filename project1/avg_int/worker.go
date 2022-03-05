package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func worker(data DataForWorker, results chan<- WorkerMessage) {

	ints := make([]int, 0)

	var message WorkerMessage
	var sum int

	fmt.Printf("Starting worker on file: %s, with start: %d, end: %d\n", data.Datafile, data.Start, data.End)

	file, err := os.Open(data.Datafile)
	if err != nil {
		log.Panic("worker crashed reading file: ", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Panic("worker crashed reading file stats: ", err)
	}

	file.Seek(int64(data.Start-1), io.SeekStart)

	buf := make([]byte, data.End-data.Start+1)
	n, err := file.Read(buf[:cap(buf)])
	buf = buf[:n]
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	str := string(buf)

	if len(str) == 0 {
		return
	}

	strIntegers := strings.Fields(str)

	for _, s := range strIntegers {
		convInt, err := strconv.Atoi(s)
		if err != nil {
			log.Panic("Error parsing int: ", err)
		}
		ints = append(ints, convInt)
	}

	var start, end int
	start = 0
	end = len(ints)

	message.Pcount = len(ints) - 2

	first := str[0]
	if unicode.IsSpace(rune(first)) || data.Start == 1 {
		message.Prefix = ""
		message.Pcount = message.Pcount + 1
	} else {
		message.Prefix = strIntegers[0]
		start = 1
	}

	last := str[len(str)-1]

	if unicode.IsSpace(rune(last)) || fileInfo.Size() == int64(data.End) {
		message.Suffix = ""
		message.Pcount = message.Pcount + 1
	} else {
		message.Suffix = strIntegers[len(strIntegers)-1] //ints[len(ints)-1]
		end = end - 1
	}

	for i := start; i < end; i++ {
		sum = sum + ints[i]
	}

	message.Start = data.Start
	message.End = data.End

	message.Psum = sum

	results <- message

}
