package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type krakenApiResponse struct {
	Result btcPriceResponse `json:"result"`
}

type btcPriceResponse struct {
	BtcUsd lastTradedPrice `json:"XXBTZUSD,omitempty"`
	BtcChf lastTradedPrice `json:"XBTCHF,omitempty"`
	BtcEur lastTradedPrice `json:"XXBTZEUR,omitempty"`
}

type lastTradedPrice struct {
	LastTradedPrice []string `json:"c"`
}

type BtcLastTradedPrice struct {
	Pair   string `json:"pair"`
	Amount string `json:"amount"`
}

func FetchLastTradedPrice(wg *sync.WaitGroup, btcPair string, ltp chan<- BtcLastTradedPrice, errorChan chan<- error) {
	defer wg.Done()

	krakenApiUrl := os.Getenv("API_LAST_TRADED_PRICE")
	if krakenApiUrl == "" {
		errorChan <- errors.New("API_LAST_TRADED_PRICE environment variable must be defined")
		return
	}

	url := fmt.Sprintf("%s?pair=%s", krakenApiUrl, btcPair)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Printf("error on creating api request %s", err)
		errorChan <- err
		return
	}
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("error on api request %s", err)
		errorChan <- err
		return
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("could not close response body: %s", err)
		}
	}()

	var lastTradedPriceResponse krakenApiResponse
	if err := json.NewDecoder(res.Body).Decode(&lastTradedPriceResponse); err != nil {
		log.Printf("error decoding api response %s", err)
		errorChan <- err
		return
	}

	switch btcPair {
	case "BTCUSD":
		ltp <- BtcLastTradedPrice{Pair: "BTC/USD", Amount: lastTradedPriceResponse.Result.BtcUsd.LastTradedPrice[0]}
	case "BTCCHF":
		ltp <- BtcLastTradedPrice{Pair: "BTC/CHF", Amount: lastTradedPriceResponse.Result.BtcChf.LastTradedPrice[0]}
	case "BTCEUR":
		ltp <- BtcLastTradedPrice{Pair: "BTC/EUR", Amount: lastTradedPriceResponse.Result.BtcEur.LastTradedPrice[0]}
	default:
		errorChan <- errors.New("no valid BTC pair given")
	}
}
