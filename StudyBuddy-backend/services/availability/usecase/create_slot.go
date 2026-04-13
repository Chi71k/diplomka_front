package usecase

import (
	"fmt"
	"studybuddy/backend/services/availability/domain"
	"time"
)

type CreateSlotInput struct {
	UserID    string
	DayOfWeek int    // 0 = Monday ... 6 = Sunday
	StartTime string // "HH:MM"
	EndTime   string // "HH:MM"
	Timezone  string
}

type CreateSlotOutput struct {
	Slot domain.Slot
}

type CreateSlot interface {
	CreateSlot(input CreateSlotInput) (CreateSlotOutput, error)
}

type createSlot struct {
	repo SlotRepository
}

func NewCreateSlot(repo SlotRepository) CreateSlot {
	return &createSlot{repo: repo}
}

func (s *createSlot) CreateSlot(input CreateSlotInput) (CreateSlotOutput, error) {
	if input.DayOfWeek < 0 || input.DayOfWeek > 6 {
		return CreateSlotOutput{}, domain.ErrInvalidDayOfWeek
	}

	tz, err := time.LoadLocation(input.Timezone)
	if err != nil {
		return CreateSlotOutput{}, domain.ErrInvalidTimezone
	}

	start, err := parseTime(input.StartTime, tz)
	if err != nil {
		return CreateSlotOutput{}, fmt.Errorf("%w: start_time %q", domain.ErrInvalidTimeFormat, input.StartTime)
	}

	end, err := parseTime(input.EndTime, tz)
	if err != nil {
		return CreateSlotOutput{}, fmt.Errorf("%w: end_time %q", domain.ErrInvalidTimeFormat, input.EndTime)
	}

	slot := &domain.Slot{
		UserID:    input.UserID,
		DayOfWeek: input.DayOfWeek,
		StartTime: start,
		EndTime:   end,
		Timezone:  input.Timezone,
	}

	if err := slot.Validate(); err != nil {
		return CreateSlotOutput{}, err
	}

	if err := s.repo.Create(slot); err != nil {
		return CreateSlotOutput{}, err
	}

	return CreateSlotOutput{Slot: *slot}, nil
}
