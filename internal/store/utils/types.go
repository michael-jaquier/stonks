package utils

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

type Results struct {
	R []Result `json:"results"`
}

type MetaData struct {
	MetaData []string `json:"Meta Data"`
}
