package main

import (
	"errors"
	"fmt"
	"strings"
)

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

	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
	OrderBy string `json:"order_by"`
}

// Struct for searching by number parameters
// Rules of matching:
// Value between Min and Max (if greater than 0)
// Or value is one of CustomVariants (if contains elements)
type NumberSearchParameters struct {
	Min            float64 `json:"min"`
	Max            float64 `json:"max"`
	CustomVariants []string   `json:"custom_variants"`
}
// Priority 1 if Min or Max greater then 0
// Priority 2 if Min or Max are 0 and CustomVariants contains elements
// Error in other cases
func (p NumberSearchParameters) getWhereCondition (columnName string) (string, error) {
	if p.Min > 0 {
		if p.Max > 0 {
			condition := fmt.Sprintf("%v between %v and %v", columnName, p.Min, p.Max)
			return condition,nil
		} else {
			condition := fmt.Sprintf("%v >= %v", columnName, p.Min)
			return condition,nil
		}
	} else {
		if p.Max > 0 {
			condition := fmt.Sprintf("%v <= %v", columnName, p.Max)
			return condition, nil
		} else {
			if len(p.CustomVariants) > 0 {
				variantsList := strings.Join(p.CustomVariants, ",")
				condition := fmt.Sprintf("%v in (%v)", columnName, variantsList)
				return condition, nil
			}
		}
	}
	errText := fmt.Sprintf("Wrong or empty conditions for number field %v", columnName)
	return "", errors.New(errText)
}

type CountStruct struct {
	Count int `json:"count"`
}

type AparmentsApiResponse struct {
	Results []Apartment `json:"results,omitempty"`
	Count   int         `json:"count,omitempty"`
}