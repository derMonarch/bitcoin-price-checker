package ltp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/derMonarch/bitcoin-price-checker/internal/client"
	"log"
	"net/http"
	"sync"
	"time"
)

type BtcLtpResponse struct {
	Ltp *[]client.BtcLastTradedPrice `json:"ltp"`
}

func LastTradedPriceHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	if !params.Has("pairs") {
		pairs, err := fetchPairs("BTCUSD", "BTCCHF", "BTCEUR")
		if err != nil {
			log.Printf("could not fetch all pairs: %s", err)
			http.Error(w, fmt.Sprintf("could not process request: %s", err), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(BtcLtpResponse{Ltp: pairs}); err != nil {
			log.Println("could not encode struct value")
		}

		return
	}

	pairs, err := fetchPairs(params["pairs"]...)
	if err != nil {
		log.Printf("could not fetch all pairs: %s", err)
		http.Error(w, fmt.Sprintf("could not process request: %s", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(BtcLtpResponse{Ltp: pairs}); err != nil {
		log.Println("could not encode struct value")
	}
}

func fetchPairs(pairs ...string) (*[]client.BtcLastTradedPrice, error) {
	ltpChan := make(chan client.BtcLastTradedPrice, len(pairs))
	errorChan := make(chan error)

	var wg sync.WaitGroup

	go func() {
		wg.Wait()
		close(ltpChan)
		close(errorChan)
	}()

	for _, p := range pairs {
		wg.Add(1)
		go client.FetchLastTradedPrice(&wg, p, ltpChan, errorChan)
	}

	btcPairs := make([]client.BtcLastTradedPrice, 0, len(pairs))

	for {
		select {
		case errMessage := <-errorChan:
			return nil, errMessage
		case ltp := <-ltpChan:
			btcPairs = append(btcPairs, ltp)
			if len(btcPairs) == len(pairs) {
				return &btcPairs, nil
			}
		case <-time.After(15 * time.Second):
			err := errors.New("timeout calling api")
			return nil, err
		default:
		}
	}
}
