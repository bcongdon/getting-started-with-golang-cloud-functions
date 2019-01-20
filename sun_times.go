package helloworld

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/schema"
	"github.com/nathan-osman/go-sunrise"
)

type Request struct {
	Lat  float64   `schema:"lat,required"`
	Lon  float64   `schema:"lon,required"`
	Date time.Time `schema:"date,required"`
}

type Response struct {
	Sunset  time.Time `json:"sunset"`
	Sunrise time.Time `json:"sunrise"`
}

func SunTimes(w http.ResponseWriter, r *http.Request) {
	var decoder = schema.NewDecoder()
	decoder.RegisterConverter(time.Time{}, dateConverter)

	// Parse query from query string
	var req Request
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	// Perform sunrise/sunset calculation
	sunrise, sunset := sunrise.SunriseSunset(
		req.Lat, req.Lon,
		req.Date.Year(), req.Date.Month(), req.Date.Day(),
	)

	// Send response back to client
	w.WriteHeader(http.StatusOK)
	response := Response{sunset, sunrise}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		panic(err)
	}
}

func dateConverter(value string) reflect.Value {
	s, _ := time.Parse("2006-01-_2", value)
	return reflect.ValueOf(s)
}

func main() {
	http.HandleFunc("/", SunTimes)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
