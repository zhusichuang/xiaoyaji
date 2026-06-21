package util

import "time"

func ParseRFC3339(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func TodayRangeByOffset(offsetMin int) (time.Time, time.Time, string) {
	nowUTC := time.Now().UTC()
	localNow := nowUTC.Add(time.Duration(-offsetMin) * time.Minute)

	localStart := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, time.UTC)
	localEnd := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 23, 59, 59, 0, time.UTC)

	startUTC := localStart.Add(time.Duration(offsetMin) * time.Minute)
	endUTC := localEnd.Add(time.Duration(offsetMin) * time.Minute)
	return startUTC, endUTC, localNow.Format("2006-01-02")
}

func LocalActionTimeFromClock(clock string, offsetMin int) time.Time {
	nowUTC := time.Now().UTC()
	localNow := nowUTC.Add(time.Duration(-offsetMin) * time.Minute)
	parsed, err := time.Parse("15:04", clock)
	if err != nil {
		return time.Now().UTC()
	}

	localValue := time.Date(
		localNow.Year(),
		localNow.Month(),
		localNow.Day(),
		parsed.Hour(),
		parsed.Minute(),
		0,
		0,
		time.UTC,
	)
	return localValue.Add(time.Duration(offsetMin) * time.Minute)
}
