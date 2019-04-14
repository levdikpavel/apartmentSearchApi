package main

type Apartment struct {
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

type ApartmentSearchRequest struct {
	City                   string `json:"city"`
	District               string `json:"district"`
	Address                string `json:"address"`
	ResidentalCompoundName string `json:"residental_compound_name"`

	CorpusName       string                 `json:"corpus_name"`
	FloorsCountRange NumberSearchParameters `json:"floors_count_range"`

	FloorRange      NumberSearchParameters `json:"floor_range"`
	RoomsCountRange NumberSearchParameters `json:"rooms_count_range"`
	SquareRange     NumberSearchParameters `json:"square_range"`
	CostRange       NumberSearchParameters `json:"cost_range"`
}

type NumberSearchParameters struct {
	Min            float64 `json:"min"`
	Max            float64 `json:"max"`
	CustomVariants []int   `json:"custom_variants"`
}