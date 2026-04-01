package gcal

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"studybuddy/backend/services/availability/domain"
)

const importWindow = 4 * 7 * 24 * time.Hour
const minSlotOccurrences = 2

func (p *Provider) ImportEvents(
	ctx context.Context,
	conn *domain.GCalConnection,
	userID string,
) ([]domain.Slot, error) {
	svc, err := p.calendarService(ctx, conn)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	events, err := svc.Events.List(conn.CalendarID).
		TimeMin(now.Format(time.RFC3339)).
		TimeMax(now.Add(importWindow).Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime").
		ShowDeleted(false).
		MaxResults(2500).
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("list calendar events: %w", err)
	}

	return eventsToSlots(events.Items, userID), nil
}

// helpers
type slotKey struct {
	dayOfWeek int
	startHHMM string // "HH:MM"
	endHHMM   string // "HH:MM"
	timezone  string // IANA name
}

func eventsToSlots(items []*calendar.Event, userID string) []domain.Slot {
	counts := make(map[slotKey]int)

	for _, item := range items {
		if item.Start == nil {
			continue
		}

		if item.Start.DateTime == "" {
			continue
		}

		key, ok := eventToSlotKey(item)
		if !ok {
			continue
		}
		counts[key]++
	}

	var slots []domain.Slot
	for key, count := range counts {
		if count < minSlotOccurrences {
			continue
		}

		loc, err := time.LoadLocation(key.timezone)
		if err != nil {
			log.Printf("gcal import: skip slot with unknown timezone %q: %v", key.timezone, err)
			continue
		}

		start, err := time.ParseInLocation("15:04", key.startHHMM, loc)
		if err != nil {
			continue
		}
		end, err := time.ParseInLocation("15:04", key.endHHMM, loc)
		if err != nil {
			continue
		}

		slots = append(slots, domain.Slot{
			UserID:    userID,
			DayOfWeek: key.dayOfWeek,
			StartTime: start,
			EndTime:   end,
			Timezone:  key.timezone,
		})
	}

	return slots
}

func eventToSlotKey(item *calendar.Event) (slotKey, bool) {
	start, err := time.Parse(time.RFC3339, item.Start.DateTime)
	if err != nil {
		return slotKey{}, false
	}
	end, err := time.Parse(time.RFC3339, item.End.DateTime)
	if err != nil {
		return slotKey{}, false
	}

	tz := item.Start.TimeZone
	if tz == "" {
		_, offset := start.Zone()
		tz = fixedOffsetName(offset)
	}

	if _, err := time.LoadLocation(tz); err != nil {
		tz = "UTC"
		start = start.UTC()
		end = end.UTC()
	}

	loc, _ := time.LoadLocation(tz)
	localStart := start.In(loc)
	localEnd := end.In(loc)

	dow := isoWeekday(localStart.Weekday())

	return slotKey{
		dayOfWeek: dow,
		startHHMM: localStart.Format("15:04"),
		endHHMM:   localEnd.Format("15:04"),
		timezone:  tz,
	}, true
}

func isoWeekday(w time.Weekday) int {
	if w == time.Sunday {
		return 6
	}
	return int(w) - 1
}

func fixedOffsetName(offsetSeconds int) string {
	hours := offsetSeconds / 3600
	if hours >= 0 {
		return fmt.Sprintf("Etc/GMT-%d", hours)
	}
	return fmt.Sprintf("Etc/GMT+%d", -hours)
}

func (p *Provider) calendarService(ctx context.Context, conn *domain.GCalConnection) (*calendar.Service, error) {
	token := &oauth2.Token{
		AccessToken:  conn.AccessToken,
		RefreshToken: conn.RefreshToken,
		Expiry:       conn.TokenExpiry,
		TokenType:    "Bearer",
	}

	httpClient := p.oauthCfg.Client(ctx, token)

	svc, err := calendar.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("create calendar service: %w", err)
	}
	return svc, nil
}
