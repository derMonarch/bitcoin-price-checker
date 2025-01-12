package internal

import (
	"encoding/json"
	"github.com/derMonarch/bitcoin-price-checker/internal/ltp"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLtpNoParams(t *testing.T) {
	t.Run("returns all pairs", func(t *testing.T) {
		t.Setenv("API_LAST_TRADED_PRICE", "https://api.kraken.com/0/public/Ticker")

		request, _ := http.NewRequest(http.MethodGet, "/api/v1/ltp", nil)
		response := httptest.NewRecorder()

		ltp.LastTradedPriceHandler(response, request)

		var btcResp ltp.BtcLtpResponse
		if err := json.NewDecoder(response.Body).Decode(&btcResp); err != nil {
			t.Fatalf("could not decode response: %s", err)
		}

		if len(*btcResp.Ltp) != 3 {
			log.Fatal("invalid response, length must be 3")
		}

		for _, l := range *btcResp.Ltp {
			switch l.Pair {
			case "BTC/EUR", "BTC/CHF", "BTC/USD":
			default:
				t.Fatal("btc pairs must be one or more of BTCEUR / BTCCHF / BTCUSD")
			}
		}
	})
}

func TestLtpEUR(t *testing.T) {
	t.Run("returns all pairs", func(t *testing.T) {
		t.Setenv("API_LAST_TRADED_PRICE", "https://api.kraken.com/0/public/Ticker")

		request, _ := http.NewRequest(http.MethodGet, "/api/v1/ltp?pairs=BTCEUR", nil)
		response := httptest.NewRecorder()

		ltp.LastTradedPriceHandler(response, request)

		var btcResp ltp.BtcLtpResponse
		if err := json.NewDecoder(response.Body).Decode(&btcResp); err != nil {
			t.Fatalf("could not decode response: %s", err)
		}

		if len(*btcResp.Ltp) != 1 {
			log.Fatal("invalid response, length must be 1")
		}

		for _, l := range *btcResp.Ltp {
			switch l.Pair {
			case "BTC/EUR", "BTC/CHF", "BTC/USD":
			default:
				t.Fatal("btc pairs must be one or more of BTCEUR / BTCCHF / BTCUSD")
			}
		}
	})
}

func TestLtpInvalidPair(t *testing.T) {
	t.Run("returns all pairs", func(t *testing.T) {
		t.Setenv("API_LAST_TRADED_PRICE", "https://api.kraken.com/0/public/Ticker")

		request, _ := http.NewRequest(http.MethodGet, "/api/v1/ltp?pairs=BTCC", nil)
		response := httptest.NewRecorder()

		ltp.LastTradedPriceHandler(response, request)

		errMessage := response.Body.String()
		if errMessage != "could not process request: no valid BTC pair given\n" {
			t.Fatalf("wrong error response, should be %s", "could not process request: no valid BTC pair given\n")
		}
	})
}
