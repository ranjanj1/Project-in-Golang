package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
)

func main() {
	// console input parameters
	// quantity of workers
	var m int
	// relative path to the text file
	var fname string

	// flags for the console parameters
	flag.IntVar(&m, "M", 1, "Specify number of workers. Default is 1")
	flag.StringVar(&fname, "fname", "input.txt", "Specify relative path tot he input file. Default is input.txt")

	// parsing the flags from console
	flag.Parse()

	fmt.Printf("Running the app with %d workers, parsing the file located at: \"%s\"\nTo change these parameters please use -M and -fname flags\n", m, fname)

	fileOptions, err := getFileOptions(fname, m)

	if err != nil {
		log.Panic("Can't open input file")
	}

	if float64(fileOptions.IntsCount) < float64(m) {
		log.Fatal(fmt.Errorf("Too much workers: %d, for total integers of: %d", m, fileOptions.IntsCount))
	}

	start := 1
	next := start + fileOptions.IntsPerWorker - 1

	jobs := make(chan WorkerMessage, m)

	for i := 0; i < m; i++ {
		if fileOptions.Remainder > 0 {
			next = next + 1
			fileOptions.Remainder = fileOptions.Remainder - 1
		}

		data := &DataForWorker{
			Datafile: fname,
			Start:    start,
			End:      next,
		}

		start = next + 1
		next = start + fileOptions.IntsPerWorker - 1

		go worker(*data, jobs)
	}

	prefixes := make(map[int]string)
	suffixes := make(map[int]string)

	var data CoordinatorData

	for a := 1; a <= m; a++ {
		message := <-jobs
		data.Sum = data.Sum + message.Psum
		data.Count = data.Count + message.Pcount
		if message.Prefix != "" {
			prefixes[message.Start] = message.Prefix
			if suf, ok := suffixes[message.Start-1]; ok {
				num, err := strconv.Atoi(fmt.Sprintf("%s%s", suf, message.Prefix))
				if err != nil {
					log.Panic("Error in working on concatenating prefix and suffix converions from string to int: ", err)
				}
				data.Sum = data.Sum + num
				data.Count = data.Count + 1
				delete(prefixes, message.Start)
				delete(suffixes, message.Start-1)
			}
		}
		if message.Suffix != "" {
			suffixes[message.End] = message.Suffix
			if pref, ok := prefixes[message.End+1]; ok {
				num, err := strconv.Atoi(fmt.Sprintf("%s%s", message.Suffix, pref))
				if err != nil {
					log.Panic("Error in working on concatenating prefix and suffix converions from string to int: ", err)
				}
				data.Sum = data.Sum + num
				data.Count = data.Count + 1
				delete(suffixes, message.End)
				delete(prefixes, message.End+1)
			}
		}
	}
	if len(prefixes) > 0 {
		for _, numStr := range prefixes {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				log.Panic("Error on parsing not matched prefixes to integers: ", err)
			}
			data.Sum = data.Sum + num
			data.Count = data.Count + 1
		}
	}
	if len(suffixes) > 0 {
		for _, numStr := range suffixes {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				log.Panic("Error on parsing not matched prefixes to integers: ", err)
			}
			data.Sum = data.Sum + num
			data.Count = data.Count + 1
		}
	}
	data.Average = float64(data.Sum) / float64(data.Count)

	str, err := json.Marshal(data)
	fmt.Println(string(str))
}
