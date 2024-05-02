package services

import (
	"e-calendar/cmd/modules/processor"
	"strings"

	ics "github.com/arran4/golang-ical"
)

type ICalendarService interface {
	CreateCourseEvents(course processor.CourseDto, semester string) []*ics.VEvent
}

func NewCalendarService(university string) ICalendarService {
	switch strings.ToLower(university) {
	case "ueh":
		return &UehCalendarService{}
	default:
		return nil
	}
}