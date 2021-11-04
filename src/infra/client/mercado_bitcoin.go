package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	uri = "https://www.mercadobitcoin.net/api"
)

var mbTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

var mbClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: mbTransport,
}

type mercadoBitcoin struct {
	http http.Client
}

func NewMercadoBitcoin() *mercadoBitcoin {
	return &mercadoBitcoin{
		http: *mbClient,
	}
}

func (mb *mercadoBitcoin) FetchDaySymmary(coin string, date time.Time) (*DaySummaryResponse, error) {

	endpoint := fmt.Sprintf("%s/%s/day-summary/%d/%d/%d", uri, coin, date.Year(), date.Month(), date.Day())
	resp, err := mb.http.Get(endpoint)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	var result DaySummaryResponse

	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &result, nil
}
