package timezones

import (
	"time"
)

type timeWithZone struct {
	Offset int `json:"offset"`
	Time time.Time `json:"time"`
}

func GenerateTimezoneList() [25]timeWithZone {
	const UTC_OFFSETS int = 25

	currentTime := time.Now()
	var timezonesList [UTC_OFFSETS]timeWithZone
	for i := range UTC_OFFSETS {
		currentOffset := UTC_OFFSETS / 2 * (-1) + i
		timezonesList[i] = timeWithZone{ currentOffset, currentTime.Add(time.Hour * time.Duration(currentOffset)) }
	}
	return timezonesList
}