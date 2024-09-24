package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/seb-emmot/code_samples/tollcalculator-go/tollcalculator"
)

type TollRequest struct {
	Passes      []string                   `json:"passes"`
	VehicleType tollcalculator.VehicleType `json:"vehicletype"`
}

type TollReponse struct {
	Fee int `json:"fee"`
}

type TollApp struct {
	tc tollcalculator.TollCalculator
}

func main() {
	mux := http.NewServeMux()

	// holidayProviderSWE can be supplemented with another supplier of holidays.
	tc, err := tollcalculator.NewTollCalculator(holidayProviderSWE)

	if err != nil {
		log.Fatal("Failed initialization")
	}

	t := TollApp{
		tc: tc,
	}

	mux.HandleFunc("/tolls", t.getTollFees)

	http.ListenAndServe(":8080", mux)
}

func holidayProviderSWE(t time.Time) bool {
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return true
	}
	if t.Year() == 2013 {
		month := t.Month()
		day := t.Day()
		if month == 1 && day == 1 ||
			month == 3 && (day == 28 || day == 29) ||
			month == 4 && (day == 1 || day == 30) ||
			month == 5 && (day == 1 || day == 8 || day == 9) ||
			month == 6 && (day == 5 || day == 6 || day == 21) ||
			month == 7 ||
			month == 11 && day == 1 ||
			month == 12 && (day == 24 || day == 25 || day == 26 || day == 31) {
			return true
		}
	}
	return false
}

func (t TollApp) getTollFees(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request TollRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Println("incoming request with", request)

	if request.VehicleType > tollcalculator.Other {
		http.Error(w, "Bad Request, invalid vehicletype in json", http.StatusBadRequest)
		return
	}

	times := []time.Time{}
	layout := "2006-01-02T15:04:05Z"
	for _, t := range request.Passes {
		parsed, err := time.Parse(layout, t)

		if err != nil {
			http.Error(w, "Bad Request, invalid timestring in json", http.StatusBadRequest)
		}

		times = append(times, parsed)
	}

	fee, err := t.tc.GetTollFees(times, request.VehicleType)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := TollReponse{Fee: fee}

	json.NewEncoder(w).Encode(response)
}
