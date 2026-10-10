package timezones

import "time"

type TimeWithZone struct {
	Offset int `json:"offset"`
	Time time.Time `json:"time"`
}

func GenerateTimezoneList() []TimeWithZone {
	const utcOffsets = 25

	currentTime := time.Now()
	timezonesList := make([]TimeWithZone, 0, utcOffsets)
	for offset := -12; offset <= 12; offset++ {
		location := time.FixedZone("UTC", offset * 3600)
		timezonesList = append(timezonesList, TimeWithZone{
			Offset: offset,
			Time: currentTime.In(location),
		})
	}

	return timezonesList
}