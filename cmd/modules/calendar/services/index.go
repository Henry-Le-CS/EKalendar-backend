package calender_services

import (
	"e-calendar/cmd/modules/processor"
	processor_srv "e-calendar/cmd/modules/processor/services"
	"log"
	"strings"

	ics "github.com/arran4/golang-ical"
)

type ICalendarService interface {
	CreateCourseEvents(course processor.CourseDto, semester string) ([]*ics.VEvent, error)
	CreateCalendar(input string) (string, error)
}

func NewCalendarService(university string, processorSrv processor_srv.IProcessorService) ICalendarService {
	if processorSrv == nil {
		log.Fatalf("Processor service at university %s is nil\n", university)
		return nil
	}

	switch strings.ToLower(university) {
	case "ueh":
		return &UehCalendarService{
			processorSv: processorSrv,
		}
	default:
		return nil
	}
}