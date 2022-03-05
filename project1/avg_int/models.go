package main

// dataForWorker will be used to spawn workers with needed parameters
type DataForWorker struct {
	Datafile string `json:"datafile"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
}

type FileOptions struct {
	ASCIICount    int `json:"ascii_count"`
	IntsPerWorker int `json:"ints_per_worker"`
	Remainder     int `json:"remainder"`
	IntsCount     int `json:"int_count"`
}

type WorkerMessage struct {
	Psum   int    `json:"psum"`
	Pcount int    `json:"pcount"`
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
	Start  int    `json:"start"`
	End    int    `json:"end"`
}

type CoordinatorData struct {
	Sum     int     `json:"sum"`
	Count   int     `json:"count"`
	Average float64 `json:"average"`
}
