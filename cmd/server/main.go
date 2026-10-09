package main

import (
	"net/http"
	"time"
	"encoding/json"
)

type timeWithZone struct {
	Offset int `json:"offset"`
	Time time.Time `json:"time"`
}

func timezones(w http.ResponseWriter, r *http.Request) {
	timezonesList := generateTimezoneList()
	jsonStr, err := json.Marshal(timezonesList)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}

func generateTimezoneList() [25]timeWithZone {
	const UTC_OFFSETS int = 25
	currentTime := time.Now() // TODO: Doesn't respect actual timezone of PC
	var timezonesList [UTC_OFFSETS]timeWithZone
	for i := range UTC_OFFSETS {
		currentOffset := UTC_OFFSETS / 2 * (-1) + i
		timezonesList[i] = timeWithZone{ currentOffset, currentTime.Add(time.Hour * time.Duration(currentOffset)) }
	}
	return timezonesList
}

func main() {
	http.HandleFunc("/timezones", timezones)
	http.ListenAndServe(":8080", nil)
}