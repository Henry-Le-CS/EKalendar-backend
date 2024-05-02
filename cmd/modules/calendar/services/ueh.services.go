package calender_services

import (
	"e-calendar/cmd/common"
	"e-calendar/cmd/modules/processor"
	"fmt"
	"strconv"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
)

type UehCalendarService struct {
}

func (service *UehCalendarService) CreateCalendar(courseList processor.CourseListDto) (string, error) {
	calendar := ics.NewCalendar()
	calendar.SetMethod(ics.MethodRequest)

	for _, course := range courseList.Courses {
		events, err := service.CreateCourseEvents(course, courseList.Semester)

		if err != nil {
			fmt.Printf("Create calendar failed by execption: %s", err)
			return "", err
		}

		for _, event := range events {
			calendar.AddVEvent(event)
		}
	}

	calendar.SetTimezoneId("Asia/Ho_Chi_Minh")
	
	res := calendar.Serialize()
	return res, nil
}

type CourseWithScheduleDto struct {
	processor.CourseInfo
	processor.CourseDetail
	processor.LecturerInfo
	processor.ScheduleDto
}

func (s *UehCalendarService) CreateCourseEvents(course processor.CourseDto, semester string) ([]*ics.VEvent, error) {
	if len(course.Schedule) == 0 {
		return nil, nil
	}

	events := make([]*ics.VEvent, 0)

	for _, schedule := range course.Schedule {
		event := &ics.VEvent{}

		dto := CourseWithScheduleDto{
			CourseInfo: course.CourseInfo,
			CourseDetail: course.CourseDetail,
			LecturerInfo: course.LecturerInfo,
			ScheduleDto: schedule,
		}

		s.SetEventTitle(event,dto)
		s.SetCourseDescription(event, dto)
		s.SetCourseLocation(event, dto)


		if err := s.ScheduleForEvent(event, dto); err != nil {
			return nil, err
		}
		
		events = append(events, event)
	}

	return events, nil
}

func (s *UehCalendarService) SetEventTitle(event *ics.VEvent, course CourseWithScheduleDto) {
	str := course.CourseName + " - " + course.ClassCode
	event.SetSummary(str)
}

func (service *UehCalendarService) SetCourseDescription(event *ics.VEvent,course CourseWithScheduleDto) {
	description := ""

	description += "Course: " + course.CourseName + "\n"
	description += "Course Code: " + course.CourseCode + "\n"
	description += "Credits: " + string(course.Credits) + "\n"
	description += "Class Code: " + course.ClassCode + "\n"

	description += "Lecturer: " + course.Lecturer + "\n"
	description += "Lecturer Email: " + course.LecturerEmail + "\n"

	description += "Room: " + course.Room + "\n"
	description += "Campus: " + course.Campus + "\n"

	event.SetDescription(description)
}

func (service *UehCalendarService) SetCourseLocation(event *ics.VEvent, course CourseWithScheduleDto) {
	event.SetLocation(course.Address)
}

func (service *UehCalendarService) ScheduleForEvent(event *ics.VEvent, course CourseWithScheduleDto) error {
	startDate, endDate, err := service.getStartEndDate(course)

	if err != nil {
		return err
	}

	startTime, endTime, err := service.GetStartEndTime(course.Session, startDate)

	if err != nil {
		fmt.Printf("Error when setting start and end time, session: %s, startDate: %s", course.Session, startDate)
		return err
	}

	event.SetStartAt(startTime)
	event.SetEndAt(endTime)

	if course.StartDate != course.EndDate {
		rule := "FREQ=WEEKLY;UNTIL=" + endDate.Format("20060102T150405Z")
		event.AddRrule(rule)
	}

	return nil
}

func (service *UehCalendarService) GetStartEndTime(session string, startTime time.Time) (time.Time, time.Time, error) {
    // Get time from session
    chunks := strings.Split(session, " (")
    if len(chunks) != 2 {
        return time.Time{}, time.Time{}, fmt.Errorf("invalid session format: %s", session)
    }

    // sessionStart := chunks[0]
    sessionStart := chunks[1][:len(chunks[1])-1]

    start, err := service.caculateStartTime(sessionStart, startTime)

	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error calculating start time: %w", err)
	}

    end, err := service.calculateEndTime(start, sessionStart)

    if err != nil {
        return time.Time{}, time.Time{}, fmt.Errorf("error calculating end time: %w", err)
    }

    return start, end, nil
}

func (service *UehCalendarService) caculateStartTime(start string, startTime time.Time) (time.Time, error) {
	// Get time from session
	chunks := strings.Split(start, "g")
	minute := 0

	// For example: 7g30
	if len(chunks) == 2 {
		minute, _ = strconv.Atoi(chunks[1])
	}

	hour, err := strconv.Atoi(chunks[0])
	
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing start time: %w", err)
	}

	startTime = time.Date(
		startTime.Year(), 
		startTime.Month(), 
		startTime.Day(), 
		hour, 
		minute, 
		0, 
		0, 
		time.Local,
	)

	return startTime, nil
}

func (service *UehCalendarService) calculateEndTime(startTime time.Time, durationStr string) (time.Time, error) {
    durationParts := strings.Split(durationStr, "g")

	minutes := 0

    hour, err := strconv.Atoi(durationParts[0])
    if err != nil {
        return time.Time{}, fmt.Errorf("error parsing duration hour: %w", err)
    }
	

	if len(durationParts) == 2 {
		minutes, err = strconv.Atoi(durationParts[1])

		if err != nil {
			return time.Time{}, fmt.Errorf("error parsing duration minutes: %w", err)
		}
	}

	endTime := time.Date(
		startTime.Year(),
		startTime.Month(),
		startTime.Day(),
		// End time should be start time + duration ( which is an hour now)
		hour + 1,
		minutes,
		0,
		0,
		time.Local,
	)

    return endTime, nil
}

func (service *UehCalendarService) countWeeksFromTimeRange(start, end time.Time) int64 {
	difference := end.Sub(start)
	weekCounts := int64(difference.Hours() / 24 / 7)

	return weekCounts
}

func (service *UehCalendarService) getStartEndDate(course CourseWithScheduleDto) (time.Time, time.Time, error) {
	startDate, err := common.ParseDate(course.StartDate,"dd/MM/yyyy")

	if err != nil {
		fmt.Printf("Error when parsing start date %s\n", course.StartDate)
		return time.Time{}, time.Time{}, err
	}


	endDate, err := common.ParseDate(course.EndDate,"dd/MM/yyyy")

	if err != nil {
		fmt.Println("Error when parsing end date")
		return time.Time{}, time.Time{}, err
	}

	return startDate, endDate, nil
}