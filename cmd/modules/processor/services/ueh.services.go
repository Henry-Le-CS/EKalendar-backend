package processor_srv

import (
	"e-calendar/cmd/modules/processor"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type UehProcessorService struct {}


func NewUehProcessorService() *UehProcessorService {
	return &UehProcessorService{}
}

func (s *UehProcessorService) ProcessFullPage(input string) (processor.CourseListDto, error) {
	courseListDto := &processor.CourseListDto{}
	lines := strings.Split(input, "\n")

	block := []string{}
	isProcessing := false
	re, err := regexp.Compile(`Năm học:\s*(\d+)\s*-\s*Học kỳ:\s*([A-Za-z]+)`)

	if err != nil {
		return *courseListDto, err
	}

	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "Mã LHP") {
			isProcessing = true
		}

		// If line match this pattern => Năm học: 2024 - Học kỳ: HKD, get year and semester
		if matches := re.FindStringSubmatch(lines[i]); len(matches) >= 3 {
			s.processSemesterYear(courseListDto, matches)
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
			course, err := s.ProcessCourse(courseBlock)

			if err != nil {
				return *courseListDto, err
			}
			
			courseListDto.AddCourse(course)

			isProcessing = false
			block = []string{}
		}

		if isAtLastCourse {
			break
		}
	}

	if !courseListDto.IsValid() {
		return *courseListDto, fmt.Errorf("không thể xử lí được danh sách khóa học")
	}
	
	return *courseListDto, nil
}

func (s *UehProcessorService) ProcessCourse(input string) (processor.CourseDto, error) {
	/*
2	Mã LHP: 24D1ECO50117901
Tên HP: Phát triển bất động sản nâng cao (ECO501179)
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
Ba	2 (07g10)	N2-311	05/03/2024->05/03/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
*/
	lines := strings.Split(input, "\n")

	if len(lines) <= 2 {
		return processor.CourseDto{}, fmt.Errorf("không thể xử lí được danh sách khóa học, số lượng dòng quá ít")
	}

	// We skip the first line
	courseInfo, err := s.processCourseName(lines[1])

	if err != nil {
		return processor.CourseDto{}, err
	}

	lecturerInfo, err := s.processLecturer(lines[2])

	if err != nil {
		return processor.CourseDto{}, err
	}

	classInfo, err := s.processClassInfo(lines[2])

	if err != nil {
		return processor.CourseDto{}, err
	}

	courseDto := &processor.CourseDto{
		CourseInfo: courseInfo,
		LecturerInfo: lecturerInfo,
		CourseDetail: classInfo,
	}

	// Process the rest of the schedule
	for i := 2; i < len(lines); i++ {
		prcsr := s.getScheduleProcessor(i)
		schedule, err := prcsr(lines[i])

		if err != nil {
			return processor.CourseDto{}, err
		}
		
		if schedule != nil {
			courseDto.AddSchedule(*schedule)
		}
		
	}

	return *courseDto, nil
}

func (s *UehProcessorService) processCourseName(input string) (processor.CourseInfo, error) {
	/*
		Tên HP: Phát triển bất động sản nâng cao (ECO501179)
	*/
	// Split the `(` character
	chunks := strings.Split(input, " (")

	if len(chunks) <= 1 || !strings.Contains(chunks[0], "Tên HP: ") {
		return processor.CourseInfo{}, fmt.Errorf("không thể xử lí được tên khóa học sau %s", input)
	}

	// Remove the `Tên HP: ` prefix from the course name
	courseName := strings.Split(chunks[0], ": ")[1]

	// Remove the `)` character from the course code
	courseCode := chunks[1][:len(chunks[1])-1]

	return processor.CourseInfo{
		CourseName: courseName,
		CourseCode: courseCode,
	}, nil
}

func (s *UehProcessorService) processLecturer(input string) (processor.LecturerInfo, error) {
/*
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
*/

	lecturerInfo := strings.Split(input, "\t")[0]

	chunks := strings.Split(lecturerInfo, " (")
	// Should get ["Giảng viên : Lê Nguyệt Trân", "Email: tranln@ueh.edu.vn)"]

	if len(chunks) <= 1 {
		return processor.LecturerInfo{}, nil
	}

	// Split the "Giảng viên : " prefix from the lecturer name
	if !strings.Contains(chunks[0], "Giảng viên : ") {
		return processor.LecturerInfo{}, fmt.Errorf("không thể xử lí được thông tin giảng viên sau %s", lecturerInfo)
	}
	
	lecturerName := strings.Split(chunks[0], " : ")[1]

	if !strings.Contains(chunks[1], "Email: ") {
		return processor.LecturerInfo{}, fmt.Errorf("không thể xử lí được thông tin email sau %s", lecturerInfo)
	}
	
	// Remove the "Email: " prefix from the lecturer email
	lecturerEmail := strings.Split(chunks[1], ": ")[1]

	// Remove the ")" character from the lecturer email
	lecturerEmail = lecturerEmail[:len(lecturerEmail)-1]

	return processor.LecturerInfo{
		Lecturer: lecturerName,
		LecturerEmail: lecturerEmail,
	}, nil
}

func (s *UehProcessorService) processClassInfo(input string) (processor.CourseDetail, error) {
/*
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
*/
	chunks := strings.Split(input, "\t")

	if len(chunks) <= 2 {
		return processor.CourseDetail{}, nil
	}

	// Credits are stored in the second chunk
	creds := chunks[1]

	// Class code is stored in the third chunk
	classCode := chunks[2]

	credits, err := strconv.Atoi(creds)

    if err != nil {
		return processor.CourseDetail{}, fmt.Errorf("không thể chuyển đổi số tín chỉ từ %s",err)
    }

	return processor.CourseDetail{
		ClassCode: classCode,
		Credits: credits,
	}, nil
}

func (s *UehProcessorService) processScheduleInClassInfo(input string) (*processor.ScheduleDto, error) {
/*
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
*/
	chunks := strings.Split(input, "\t")

	if len(chunks) <= 2 {
		return nil, nil
	}
	
	chunksWithSchedule := strings.Join(chunks[3:], "\t")

	return s.processSchedule(chunksWithSchedule)
}

func (s *UehProcessorService) processSchedule(input string) (*processor.ScheduleDto, error) {
	/*
		Ba	2 (07g10)	N2-311	05/03/2024->05/03/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
	*/
	
	chunks := strings.Split(input, "\t")
	
	if len(chunks) <= 4 {
		return nil, nil
	}

	for i := 0; i < len(chunks); i++ {
		if chunks[i] == "" {
			return nil, nil
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
	}, nil
}


func (s *UehProcessorService) getScheduleProcessor(line int) func (string) (*processor.ScheduleDto, error){
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