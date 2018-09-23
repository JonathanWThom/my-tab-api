package main

import "time"

func stringsToTimes(strings []string) ([]time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	realTimes := []time.Time{}

	for _, stringTime := range strings {
		convTime, err := time.Parse(layout, stringTime)
		if err != nil {
			return nil, err
		}
		realTimes = append(realTimes, convTime)
	}

	return realTimes, nil
}
