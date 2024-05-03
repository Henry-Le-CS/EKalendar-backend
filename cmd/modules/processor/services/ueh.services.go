package processor_srv

import (
	"e-calendar/cmd/modules/processor"
	"regexp"
	"strconv"
	"strings"
)

type UehProcessorService struct {}


func NewUehProcessorService() *UehProcessorService {
	return &UehProcessorService{}
}

func (s *UehProcessorService) ProcessFullPage(input string) (processor.CourseListDto, error) {
	CourseListDto := &processor.CourseListDto{}
	lines := strings.Split(input, "\n")

	block := []string{}
	isProcessing := false
	re, err := regexp.Compile(`Năm học:\s*(\d+)\s*-\s*Học kỳ:\s*([A-Za-z]+)`)

	if err != nil {
		return *CourseListDto, err
	}

	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "Mã LHP") {
			isProcessing = true
		}

		// If line match this pattern => Năm học: 2024 - Học kỳ: HKD, get year and semester
		if matches := re.FindStringSubmatch(lines[i]); len(matches) >= 3 {
			s.processSemesterYear(CourseListDto, matches)
		}
		
		if isProcessing {
			block = append(block, lines[i])
		}

		// At the end of a course, we process the block
		isAtCourseEnd := i < len(lines) - 1 && strings.Contains(lines[i + 1], "Mã LHP") && isProcessing;
		isAtLastCourse := i < len(lines) -1 && strings.Contains(lines[i + 1], "Copyright");
		isAtEOF := i == len(lines) - 1;

		shouldProcess := isAtCourseEnd || isAtLastCourse || isAtEOF
		
		if shouldProcess {
			courseBlock := strings.Join(block, "\n")
			course := s.ProcessCourse(courseBlock)
			
			CourseListDto.AddCourse(course)

			isProcessing = false
			block = []string{}
		}

		if isAtLastCourse {
			break
		}
	}
	
	return *CourseListDto, nil
}

func (s *UehProcessorService) ProcessCourse(input string) processor.CourseDto {
	/*
2	Mã LHP: 24D1ECO50117901
Tên HP: Phát triển bất động sản nâng cao (ECO501179)
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
Ba	2 (07g10)	N2-311	05/03/2024->05/03/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
*/
	lines := strings.Split(input, "\n")

	if len(lines) <= 2 {
		return processor.CourseDto{}
	}

	// We skip the first line
	courseInfo := s.processCourseName(lines[1])
	lecturerInfo := s.processLecturer(lines[2])
	classInfo := s.processClassInfo(lines[2])

	courseDto := &processor.CourseDto{
		CourseInfo: courseInfo,
		LecturerInfo: lecturerInfo,
		CourseDetail: classInfo,
	}

	// Process the rest of the schedule
	for i := 2; i < len(lines); i++ {
		prcsr := s.getScheduleProcessor(i)

		if schedule := prcsr(lines[i]); schedule != nil {
			courseDto.AddSchedule(*schedule)
		}
	}

	return *courseDto
}

func (s *UehProcessorService) processCourseName(input string) processor.CourseInfo {
	/*
		Tên HP: Phát triển bất động sản nâng cao (ECO501179)
	*/
	// Split the `(` character
	chunks := strings.Split(input, " (")

	// Remove the `Tên HP: ` prefix from the course name
	courseName := strings.Split(chunks[0], ": ")[1]

	// Remove the `)` character from the course code
	courseCode := chunks[1][:len(chunks[1])-1]

	return processor.CourseInfo{
		CourseName: courseName,
		CourseCode: courseCode,
	}
}

func (s *UehProcessorService) processLecturer(input string) processor.LecturerInfo {
/*
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
*/

	lecturerInfo := strings.Split(input, "\t")[0]

	chunks := strings.Split(lecturerInfo, " (")

	if len(chunks) <= 1 {
		return processor.LecturerInfo{
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

	return processor.LecturerInfo{
		Lecturer: lecturerName,
		LecturerEmail: lecturerEmail,
	}
}

func (s *UehProcessorService) processClassInfo(input string) processor.CourseDetail {
/*
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
*/
	chunks := strings.Split(input, "\t")

	if len(chunks) <= 2 {
		return processor.CourseDetail{
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

	return processor.CourseDetail{
		ClassCode: classCode,
		Credits: credits,
	}
}

func (s *UehProcessorService) processScheduleInClassInfo(input string) *processor.ScheduleDto {
/*
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
*/
	chunks := strings.Split(input, "\t")

	if len(chunks) <= 2 {
		return nil
	}
	
	chunksWithSchedule := strings.Join(chunks[3:], "\t")

	return s.processSchedule(chunksWithSchedule)
}

func (s *UehProcessorService) processSchedule(input string) *processor.ScheduleDto {
	/*
		Ba	2 (07g10)	N2-311	05/03/2024->05/03/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
	*/
	
	chunks := strings.Split(input, "\t")
	
	if len(chunks) <= 4 {
		return nil
	}

	for i := 0; i < len(chunks); i++ {
		if chunks[i] == "" {
			return nil
		}
	}

	day := chunks[0]
	session := chunks[1]
	room := chunks[2]
	startDate := chunks[3][:10]
	endDate := chunks[3][12:22]
	campus := chunks[4]
	address := chunks[5]

	return &processor.ScheduleDto{
		Day: day,
		Session: session,
		Room: room,
		StartDate: startDate,
		EndDate: endDate,
		Campus: campus,
		Address: address,
	}
}


func (s *UehProcessorService) getScheduleProcessor(line int) func (string) *processor.ScheduleDto{
	// Line 2 also has schedule information
	// For example: Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
	if line == 2 {
		return s.processScheduleInClassInfo
	}
	
	return s.processSchedule
}

func (s *UehProcessorService) processSemesterYear(CourseListDto *processor.CourseListDto,matches []string) {
	CourseListDto.Year = matches[1]
	CourseListDto.Semester = matches[2]
}