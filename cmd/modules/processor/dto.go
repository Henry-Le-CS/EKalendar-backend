package processor

type CalendarRequestDto struct {
	Text string
	Semester string
}

type ScheduleDto struct {
	Day string `json:"day"`
	Session string `json:"session"`
	Room string `json:"room"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Campus string `json:"campus"`
	Address string `json:"address"`
}

func (sc ScheduleDto) String() string {
	if sc.Day == "" {
		return ""
	}
	
	return `{ Day: ` + sc.Day + `, Session: ` + sc.Session + `, Room: ` + sc.Room + `, StartDate: ` + sc.StartDate + `, EndDate: ` + sc.EndDate + `, Campus: ` + sc.Campus + `, Address: ` + sc.Address + `}`
}

type CourseInfo struct {
	CourseName string `json:"course_name"`
	CourseCode string `json:"course_code"`
}

type LecturerInfo struct {
	Lecturer string `json:"lecturer"`
	LecturerEmail string `json:"lecturer_email"`
}

type CourseDetail struct {
	ClassCode string `json:"class_code"`
	Credits int `json:"credits"`
}

type CourseDto struct {
	CourseInfo
	LecturerInfo
	CourseDetail
	Schedule []ScheduleDto `json:"schedule"`
}

func (c CourseDto) String() string {
	scheduleStr := ""

	for index, schedule := range c.Schedule {
		if schedule.String() == "" {
			continue
		}

		scheduleStr += schedule.String()

		if index != len(c.Schedule) - 1 {
			scheduleStr += ", "
		}
	}

	ret := "{ "

	if c.CourseName != "" {
		ret += "CourseName: " + c.CourseName + ", "
	}

	if c.CourseCode != "" {
		ret += "CourseCode: " + c.CourseCode + ", "
	}

	if c.Lecturer != "" {
		ret += "Lecturer: " + c.Lecturer + ", "
	}

	if c.LecturerEmail != "" {
		ret += "LecturerEmail: " + c.LecturerEmail + ", "
	}

	if c.ClassCode != "" {
		ret += "ClassCode: " + c.ClassCode + ", "
	}

	if c.Credits != 0 {
		ret += "Credits: " + string(c.Credits) + ", "
	}

	if scheduleStr != "" {
		ret += "Schedule: { " + scheduleStr + " }"
	}

	// If end with ", " then remove it
	if ret[len(ret) - 2:] == ", " {
		ret = ret[:len(ret) - 2]
	}

	ret += " }"

	return ret
}

func (c *CourseDto) AddSchedule(schedule ScheduleDto) {
	if c.Schedule == nil {
		c.Schedule = make([]ScheduleDto, 0)
	}

	c.Schedule = append(c.Schedule, schedule)
}