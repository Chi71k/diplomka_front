package domain

import (
	"errors"
	"time"
)

var (
	ErrSlotNotFound      = errors.New("slot not found")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidTimeRange  = errors.New("end_time must be after start_time")
	ErrInvalidDayOfWeek  = errors.New("day_of_week must be between 0 and 6")
	ErrInvalidTimezone   = errors.New("timezone must be a valid IANA timezone name")
	ErrInvalidTimeFormat = errors.New("time must be in HH:MM 24-hour format")
	ErrSlotConflict      = errors.New("slot overlaps with existing slot")
)

type Slot struct {
	ID        string
	UserID    string
	DayOfWeek int // Example: 0=Monday ... 6=Sunday
	StartTime time.Time
	EndTime   time.Time
	Timezone  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Slot) Validate() error {
	if !s.EndTime.After(s.StartTime) {
		return ErrInvalidTimeRange
	}
	return nil
}
