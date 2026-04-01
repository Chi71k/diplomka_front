package usecase

import "time"

func parseTime(timeStr string, loc *time.Location) (time.Time, error) {
	t, err := time.ParseInLocation("15:04", timeStr, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
