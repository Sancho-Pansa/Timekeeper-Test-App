package main

import (
	"encoding/json"
	"net/http"
	"github.com/Sancho-Pansa/Timekeeper/internal/timezones"
)

func returnTimezones(w http.ResponseWriter, r *http.Request) {
	timezonesList := timezones.GenerateTimezoneList()
	jsonStr, err := json.Marshal(timezonesList)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonStr)
}

func main() {
	http.HandleFunc("/timezones", returnTimezones)
	http.ListenAndServe(":8080", nil)
}