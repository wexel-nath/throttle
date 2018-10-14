package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	address = "localhost:23456"
)

var (
	timeChan = make(chan time.Time, 10)
	timestamps []time.Time

	maxRequestRate = float64(10.0)
)

func main() {
	go recordTimestamps()

	startServer()
}

func startServer() {
	fmt.Println("Listening on: " + address)
	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/default", requestHandler)
	log.Fatal(http.ListenAndServe(address, nil))
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	rateString := r.URL.Query().Get("max")

	rate, err := strconv.ParseFloat(rateString, 64)
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		return
	}
	fmt.Printf("Setting max request rate to: %.2f requests per sec\n", rate)
	maxRequestRate = rate
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	timeChan<- time.Now()

	rate := getRequestRate()
	fmt.Printf("Request Rate: %.3f per sec\n", rate)

	if maxRequestRate == 0.0 {
		w.WriteHeader(http.StatusInternalServerError)
	} else if rate > maxRequestRate {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Write([]byte("{}"))
}

func recordTimestamps() {
	for timestamp := range timeChan {
		timestamps = append([]time.Time{timestamp}, timestamps...)
		if len(timestamps) >= 20 {
			timestamps = timestamps[:20]
		}
	}
}

func getRequestRate() float64 {
	numTimestamps := len(timestamps)
	if numTimestamps <= 1 {
		return 0.0
	}
	earliestTimestamp := timestamps[0]
	latestTimestamp := timestamps[numTimestamps - 1]

	var duration time.Duration
	if latestTimestamp.After(earliestTimestamp) {
		duration = latestTimestamp.Sub(earliestTimestamp)
	} else {
		duration = earliestTimestamp.Sub(latestTimestamp)
	}

	return float64(numTimestamps) / float64(duration.Nanoseconds()) * float64(time.Second)
}
