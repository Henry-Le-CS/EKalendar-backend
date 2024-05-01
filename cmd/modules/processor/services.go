package processor

import (
	"strconv"
	"strings"
)

type ProcessorService struct {
}

func NewProcessorService() *ProcessorService {
	return &ProcessorService{}
}

func (s *ProcessorService) ProcessCourse(input string) CourseDto {
	/*
2	Mã LHP: 24D1ECO50117901
Tên HP: Phát triển bất động sản nâng cao (ECO501179)
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
Ba	2 (07g10)	N2-311	05/03/2024->05/03/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
*/
	lines := strings.Split(input, "\n")

	// We skip the first line
	courseInfo := s.ProcessCourseName(lines[1])
	lecturerInfo := s.ProcessLecturer(lines[2])
	classInfo := s.ProcessClassInfo(lines[2])

	courseDto := &CourseDto{
		CourseInfo: courseInfo,
		LecturerInfo: lecturerInfo,
		CourseDetail: classInfo,
	}

	// Line 2 also has schedule information
	courseDto.AddSchedule(s.ProcessScheduleInClassInfo(lines[2]))

	// Process the rest of the schedule
	for i := 3; i < len(lines); i++ {
		courseDto.AddSchedule(s.ProcessSchedule(lines[i]))
	}

	return *courseDto
}

func (s *ProcessorService) ProcessCourseName(input string) CourseInfo {
	/*
		Tên HP: Phát triển bất động sản nâng cao (ECO501179)
	*/
	// Split the `(` character
	chunks := strings.Split(input, " (")

	// Remove the `Tên HP: ` prefix from the course name
	courseName := strings.Split(chunks[0], ": ")[1]

	// Remove the `)` character from the course code
	courseCode := chunks[1][:len(chunks[1])-1]

	return CourseInfo{
		CourseName: courseName,
		CourseCode: courseCode,
	}
}

func (s *ProcessorService) ProcessLecturer(input string) LecturerInfo {
/*
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
*/

	lecturerInfo := strings.Split(input, "\t")[0]

	chunks := strings.Split(lecturerInfo, " (")

	if len(chunks) <= 1 {
		return LecturerInfo{
			Lecturer: "",
			LecturerEmail: "",
		}
	}

	// Split the "Giảng viên : " prefix from the lecturer name
	lecturerName := strings.Split(chunks[0], " : ")[1]

	// Remove the "Email: " prefix from the lecturer email
	lecturerEmail := strings.Split(chunks[1], ": ")[1]

	// Remove the ")" character from the lecturer email
	lecturerEmail = lecturerEmail[:len(lecturerEmail)-1]

	return LecturerInfo{
		Lecturer: lecturerName,
		LecturerEmail: lecturerEmail,
	}
}

func (s *ProcessorService) ProcessClassInfo(input string) CourseDetail {
/*
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
*/
	chunks := strings.Split(input, "\t")

	if len(chunks) <= 2 {
		return CourseDetail{
			ClassCode: "",
			Credits: 0,
		}
	}

	// Credits are stored in the second chunk
	creds := chunks[1]

	// Class code is stored in the third chunk
	classCode := chunks[2]

	credits, err := strconv.Atoi(creds)

    if err != nil {
		panic(err)
    }

	return CourseDetail{
		ClassCode: classCode,
		Credits: credits,
	}
}

func (s *ProcessorService) ProcessScheduleInClassInfo(input string) ScheduleDto {
/*
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
*/
	chunks := strings.Split(input, "\t")

	if len(chunks) <= 2 {
		return ScheduleDto{}
	}
	
	chunksWithSchedule := strings.Join(chunks[3:], "\t")

	return s.ProcessSchedule(chunksWithSchedule)
}

func (s *ProcessorService) ProcessSchedule(input string) ScheduleDto {
	/*
		Ba	2 (07g10)	N2-311	05/03/2024->05/03/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
	*/
	

	chunks := strings.Split(input, "\t")

	if len(chunks) <= 4 {
		return ScheduleDto{}
	}

	day := chunks[0]
	session := chunks[1]
	room := chunks[2]
	startDate := chunks[3][:10]
	endDate := chunks[3][13:23]
	campus := chunks[4]
	// address := chunks[5][:len(chunks[5])-1]
	address := chunks[5]

	return ScheduleDto{
		Day: day,
		Session: session,
		Room: room,
		StartDate: startDate,
		EndDate: endDate,
		Campus: campus,
		Address: address,
	}
}