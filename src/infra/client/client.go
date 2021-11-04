package client

import "time"

type DaySummaryResponse struct {
	AvgPrice float64 `json:"avg_price"`
	Date     string  `json:"date"`
}

type IClient interface {
	FetchDaySymmary(coin string, date time.Time) (*DaySummaryResponse, error)
}
