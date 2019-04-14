package main

type Appartment struct {
	City                   string `json:"city"`
	District               string `json:"district"`
	Address                string `json:"address"`
	ResidentalCompoundName string `json:"residental_compound_name"`

	CorpusName  string `json:"corpus_name"`
	FloorsCount int    `json:"floors_count"`

	Floor      int     `json:"floor"`
	RoomsCount int     `json:"rooms_count"`
	Square     float64 `json:"square"`
	Cost       float64 `json:"cost"`
}
