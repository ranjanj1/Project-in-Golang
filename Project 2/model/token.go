package model

type Token struct {
	Id         int64
	Name       string
	Low        uint64
	Mid        uint64
	High       uint64
	PartialVal uint64
	FinalVal   uint64
}
