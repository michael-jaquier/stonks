package main

import (
	"context"
	"encoding/json"
	"fmt"
	store "github.com/michael-jaquier/stonks/internal/store/sqlite"
	"github.com/michael-jaquier/stonks/internal/store/utils"
	"os"
	"time"
)

func main() {
	file, err := os.Open("input.json")
	db, err := utils.NewDBConnection()
	if err != nil {
		panic(err)
	}
	if err := utils.DBSchemaInit(db); err != nil {
		panic(err)
	}
	s := store.New(db)

	dec := json.NewDecoder(file)
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}

		if tok == "results" {
			dec.Token() // [
			for dec.More() {
				var r = utils.Result{}
				dec.DisallowUnknownFields()
				dec.Decode(&r)
				symbolid, err := s.CreateOrGetTicker(context.Background(), r.Ticker)
				t := time.UnixMilli(int64(r.Timestamp)).UTC()
				tradingDay := time.Date(
					t.Year(), t.Month(), t.Day(),
					0, 0, 0, 0,
					time.UTC,
				)
				dailyParams := store.InsertDailyPriceParams{
					Symbolid:   symbolid.Symbolid,
					TradingDay: tradingDay,
					Open:       r.Open,
					Close:      r.Close,
				}
				insert, err := s.InsertDailyPrice(context.Background(), dailyParams)

				if err != nil {
					panic(err)
				}

				fmt.Printf("Inserts: %v\n", insert)

			}
			dec.Token() // ]
		}
	}
}
