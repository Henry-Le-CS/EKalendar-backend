package calender_services

import (
	"e-calendar/cmd/common"
	"e-calendar/cmd/modules/processor"
	processor_srv "e-calendar/cmd/modules/processor/services"
	"fmt"
	"strconv"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
)

type UehCalendarService struct {
	processorSv processor_srv.IProcessorService
}

func (s *UehCalendarService) CreateCalendar(input string) (string, error) {
	courseList, err := s.processorSv.ProcessFullPage(input)

	if err != nil {
		return "", err
	}

	calendar := ics.NewCalendar()
	calendar.SetMethod(ics.MethodRequest)

	for _, course := range courseList.Courses {
		events, err := s.CreateCourseEvents(course, courseList.Semester)

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

	description += "Môn học: " + course.CourseName + "\n"
	description += "Mã môn học: " + course.CourseCode + "\n"

	description += "Số TC: " + strconv.Itoa(course.Credits) + "\n"
	description += "Mã lớp: " + course.ClassCode + "\n"

	description += "Giảng viên: " + course.Lecturer + "\n"
	description += "Email: " + course.LecturerEmail + "\n"

	description += "Phòng: " + course.Room + "\n"
	description += "Cơ sở: " + course.Campus + "\n"

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

	tz, _ := common.TimeIn("Vietnam")

	if course.StartDate != course.EndDate {
		rule := "FREQ=WEEKLY;UNTIL=" + endDate.In(tz).Format("20060102T150405Z")
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

    start, err := service.calculateTime(sessionStart, startTime, false)

	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error calculating start time: %w", err)
	}

    end, err := service.calculateTime(sessionStart, startTime, true)

    if err != nil {
        return time.Time{}, time.Time{}, fmt.Errorf("error calculating end time: %w", err)
    }

    return start, end, nil
}

func (service *UehCalendarService) calculateTime(start string, startTime time.Time, needAdd bool) (time.Time, error) {
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

	if needAdd {
		hour += 1
	}

	tz, _ := common.TimeIn("Vietnam")
	
	startTime = time.Date(
		startTime.Year(), 
		startTime.Month(), 
		startTime.Day(), 
		hour, 
		minute, 
		0, 
		0, 
		tz,
	)

	return startTime, nil
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