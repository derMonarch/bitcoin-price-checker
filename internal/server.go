package internal

import (
	"context"
	"github.com/derMonarch/bitcoin-price-checker/internal/ltp"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func InitServer() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	wait := time.Second * 15

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/ltp", ltp.LastTradedPriceHandler).Methods("GET")

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	//blocking until signal received
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("Could not shutdown service gracefully")
	}

	log.Println("shutting down")
	os.Exit(0)
}
