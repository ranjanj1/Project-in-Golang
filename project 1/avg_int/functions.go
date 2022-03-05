package main

import (
	"io/ioutil"
	"strings"
)

func getFileOptions(path string, workers int) (*FileOptions, error) {
	// counter of ASCII integers
	count := 0
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	count = len(data)

	options := &FileOptions{
		ASCIICount:    count,
		IntsPerWorker: count / workers,
		Remainder:     count % workers,
		IntsCount:     len(strings.Fields(string(data))),
	}

	return options, nil

}
