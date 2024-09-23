package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/seb-emmot/code_samples/tollcalculator-go/tollcalculator"
)

type GetTollRequest struct {
	Passes      []string                   `json: "passes"`
	VehicleType tollcalculator.VehicleType `json: "vehicletype`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/tolls", getTollFees)

	http.ListenAndServe(":8080", mux)
}

func getTollFees(w http.ResponseWriter, r *http.Request) {
	var request GetTollRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	log.Println("incoming request with", request)

	times := []time.Time{}
	layout := "2006-01-02T15:04:05Z"
	for _, t := range request.Passes {
		parsed, err := time.Parse(layout, t)

		if err != nil {
			http.Error(w, "Bad Request, invalid timestring in json", http.StatusBadRequest)
		}

		times = append(times, parsed)
	}

	fees, err := tollcalculator.GetTollFees(times, request.VehicleType)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(fees)
}
