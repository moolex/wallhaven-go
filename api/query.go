package api

import (
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func NewQuery(q string) *QueryCond {
	return &QueryCond{
		Query:    q,
		Sorting:  "date_added",
		Order:    "desc",
		TopRange: "1M",
		Page:     1,
	}
}

func (qc *QueryCond) SetCategory(q string) *QueryCond {
	qc.Categories = q
	return qc
}

func (qc *QueryCond) SetPurity(q string) *QueryCond {
	qc.Purity = q
	return qc
}

func (qc *QueryCond) SortBy(q string) *QueryCond {
	qc.Sorting = q
	return qc
}

func (qc *QueryCond) TopList(timeRange time.Duration) *QueryCond {
	qc.SortBy(SortTopList)

	keys := []string{"1d", "3d", "1w", "1M", "3M", "6M", "1y"}
	durations := map[string]time.Duration{
		"1d": Range1day,
		"3d": Range3day,
		"1w": Range1week,
		"1M": Range1month,
		"3M": Range3months,
		"6M": Range6months,
		"1y": Range1year,
	}

	for _, k := range keys {
		if timeRange <= durations[k] {
			qc.TopRange = k
			break
		}
	}

	if qc.TopRange == "" {
		qc.TopRange = "1y"
	}

	return qc
}

func (qc *QueryCond) Asc() *QueryCond {
	qc.Order = "asc"
	return qc
}
func (qc *QueryCond) Desc() *QueryCond {
	qc.Order = "desc"
	return qc
}

func (qc *QueryCond) ToMap() (map[string]string, error) {
	if err := validate.Struct(qc); err != nil {
		return nil, err
	}

	m := map[string]string{
		"q":           qc.Query,
		"categories":  qc.Categories,
		"purity":      qc.Purity,
		"sorting":     qc.Sorting,
		"order":       qc.Order,
		"topRange":    qc.TopRange,
		"atleast":     qc.AtLeast,
		"resolutions": qc.Resolutions,
		"ratios":      qc.Ratios,
		"colors":      qc.Colors,
		"page":        strconv.Itoa(qc.Page),
	}

	for k, v := range m {
		if v == "" {
			delete(m, k)
		}
	}

	return m, nil
}
