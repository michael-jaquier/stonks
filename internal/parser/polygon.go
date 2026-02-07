package parser

import (
	"encoding/json"
	"io"
	"git"
)

type Result struct {
	Ticker         string  `json:"T"`
	Close          float64 `json:"c"`
	High           float64 `json:"h"`
	Low            float64 `json:"l"`
	Transactions   int     `json:"n"`
	Open           float64 `json:"o"`
	Timestamp      int     `json:"t"`
	Volume         int     `json:"v"`
	WeightedVolume float64 `json:"vw"`
}

func NewPolygonParser() *StreamParser[Result] {
	c := make(chan (Result), 100)
	return NewStreamProcessor(
		PolygonParser{c},
		c,
	)
}

type PolygonParser struct {
	c chan Result
}

func (p PolygonParser) Parse(r io.Reader) error {
	dec := json.NewDecoder(r)
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}

		if tok == "results" {
			dec.Token() // [
			for dec.More() {
				var r = Result{}
				dec.DisallowUnknownFields()
				dec.Decode(&r)
				var insertDailyPrice := store.Inser
				p.c <- r
			}
			dec.Token() // ]
		}
	}
	return nil

}
