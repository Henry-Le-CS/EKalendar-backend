package tests

import (
	"e-calendar/cmd/modules/processor"
	"fmt"
	"testing"
)

func TestProcessCourse(t *testing.T) {
	processorService := processor.NewProcessorService()

	result := processorService.ProcessCourse(
`2	Mã LHP: 24D1ECO50117901
Tên HP: Phát triển bất động sản nâng cao (ECO501179)
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
Ba	2 (07g10)	N2-311	05/03/2024->05/03/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM`)
	
	exp := `{ CourseName: Phát triển bất động sản nâng cao, CourseCode: ECO501179, Lecturer: Lê Nguyệt Trân, LecturerEmail: tranln@ueh.edu.vn, ClassCode: HPTC.I.PTBDS.RE.2, Credits: , Schedule: { { Day: Hai, Session: 8 (12g45), Room: N2-308, StartDate: 18/03/2024, EndDate: 3/05/2024), Campus: UEH Nguyễn Văn Linh - N2, Address: Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM}, { Day: Ba, Session: 2 (07g10), Room: N2-311, StartDate: 05/03/2024, EndDate: 5/03/2024), Campus: UEH Nguyễn Văn Linh - N2, Address: Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM} } }`
	if result.String() != exp {
		t.Fatalf("Expected %s but got %s", exp, result.String())
	}
}

func TestProcessEmptyCourse(t *testing.T) {
	processorService := processor.NewProcessorService()

	result := processorService.ProcessCourse(
`1	Mã LHP: 24D1INS53600104
Tên HP: Phí Bảo hiểm Y tế và Tai nạn (INS536001)
Giảng viên :	0`)
		
	exp := `{ CourseName: Phí Bảo hiểm Y tế và Tai nạn, CourseCode: INS536001 }`

	fmt.Println(result.String() == exp)
	fmt.Println(exp)

	if result.String() != exp {
		t.Fatalf("Expected %s but got %s", exp, result.String())
	}
}
	