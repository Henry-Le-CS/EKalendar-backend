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

	exp := `{ CourseName: Phát triển bất động sản nâng cao, CourseCode: ECO501179, Lecturer: Lê Nguyệt Trân, LecturerEmail: tranln@ueh.edu.vn, ClassCode: HPTC.I.PTBDS.RE.2, Credits: , Schedule: { { Day: Hai, Session: 8 (12g45), Room: N2-308, StartDate: 18/03/2024, EndDate: 3/05/2024, Campus: UEH Nguyễn Văn Linh - N2, Address: Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM}, { Day: Ba, Session: 2 (07g10), Room: N2-311, StartDate: 05/03/2024, EndDate: 5/03/2024, Campus: UEH Nguyễn Văn Linh - N2, Address: Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM} } }`
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

func TestFullPage(t *testing.T) {
	processorService := processor.NewProcessorService()

	result := processorService.ProcessFullPage(
`
Trang chủ
Ngành - Chương trình đào tạo
Tra cứu văn bằng
Tra cứu học phần
Hỗ trợ
31211025817 | Nguyễn Thế Thịnh
 CHỨC NĂNG
 Trang cá nhân
Thông tin cá nhân
Thông báo (11)
 Tra cứu thông tin
Chương trình đào tạo
Lịch học
Lịch thi
Quyết định sinh viên
Chuyên cần
Kết quả rèn luyện
Kết quả học tập
Tài chính sinh viên
Chi tiết hóa đơn
Xem kết quả đăng ký học phần
Học bổng, Chính sách, Miễn giảm, Trợ cấp
Học phần tương đương
 Chức năng trực tuyến
Đăng ký trường đối tác
Kết quả đăng ký vắng thi
Đăng ký tham dự lễ tốt nghiệp
Đăng ký chuyên ngành
Khảo sát nhu cầu học môn tự chọn
Khảo sát nhu cầu học
Đăng ký song ngành
Ghi danh ngành cao học
Đăng ký học phần
Đăng ký xét tốt nghiệp
Ý kiến - Thảo Luận
Đánh giá điểm rèn luyện
Ngoại trú sinh viên
Liên hệ - góp ý
Nộp chứng chỉ
Tra cứu lịch thi
Học bổng
Giấy xác nhận - bảng điểm điện tử
Thời khóa biểu
Năm học:
2024
Học kỳ:
Học kỳ đầu

 TKB Tuần
 TKB Thứ - Tiết
Năm học: 2024 - Học kỳ: HKD
STT	Mã lớp học phần
Tên học phần
CBGD	STC	Mã lớp	Thứ	Tiết (giờ) bắt đầu	Phòng	Tuần	Cơ sở	Địa chỉ
1	Mã LHP: 24D1INS53600104
Tên HP: Phí Bảo hiểm Y tế và Tai nạn (INS536001)
Giảng viên :	0
2	Mã LHP: 24D1ECO50117901
Tên HP: Phát triển bất động sản nâng cao (ECO501179)
Giảng viên : Lê Nguyệt Trân (Email: tranln@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Hai	8 (12g45)	N2-308	18/03/2024->13/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
Ba	2 (07g10)	N2-311	05/03/2024->05/03/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
3	Mã LHP: 24D1BUS53300291
Tên HP: Khởi nghiệp kinh doanh (BUS533002)
Giảng viên : Cao Văn Tiến (Email: vantien@ueh.edu.vn)	1	C.KNKD.A116_007	Hai	13 (17g45)	A116	08/04/2024->15/04/2024)	A	59C Nguyễn Đình Chiểu, P.Võ Thị Sáu, Q.3, TP.HCM
Tư	13 (17g45)	A116	10/04/2024->10/04/2024)	A	59C Nguyễn Đình Chiểu, P.Võ Thị Sáu, Q.3, TP.HCM
Sáu	13 (17g45)	A116	12/04/2024->12/04/2024)	A	59C Nguyễn Đình Chiểu, P.Võ Thị Sáu, Q.3, TP.HCM
4	Mã LHP: 24D1ADM53502488
Tên HP: Sinh hoạt lớp hk6.1[7] (ADM535024)
Giảng viên : Nguyễn Lưu Bảo Đoan (Email: doannlb@ueh.edu.vn)	0	RE002	Ba	2 (07g10)	ONLINE	14/05/2024->14/05/2024)	Online	Học online
5	Mã LHP: 24D1ARC51204701
Tên HP: Hệ thống thông tin địa lý trong quy hoạch và quản lý đô thị (ARC512047)
Giảng viên : Nguyễn Minh Quang (Email: quangnm@ueh.edu.vn), Phạm Nguyễn Hoài (Email: hoaipm@ueh.edu.vn)	3	HPTC.II.PTBDS.RE.2	Ba	2 (07g10)	N1-403	23/01/2024->23/01/2024)	UEH Nguyễn Văn Linh - N1	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
Tư	2 (07g10)	N2-210	17/01/2024->06/03/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
Tư	2 (07g10)	NGHI	03/01/2024->10/01/2024)	Online	Học online
Tư	2 (07g10)	ONLINE	31/01/2024->21/02/2024)	Online	Học online
6	Mã LHP: 24D1ECO50118001
Tên HP: Doanh nghiệp bất động sản (ECO501180)
Giảng viên : Nguyễn Thị Hồng Thu (Email: thunguyen@ueh.edu.vn), Nguyễn Thị Hoàng Thùy (Email: thuynth@ueh.edu.vn)	3	HPTC.I.PTBDS.RE.2	Ba	8 (12g45)	B2-202	12/03/2024->07/05/2024)	B2	279 Nguyễn Tri Phương P.5 Q.10 TP.HCM (Khu B2)
Ba	8 (12g45)	LMS	14/05/2024->14/05/2024)	LMS	Học trên hệ thống LMS
7	Mã LHP: 24D1ECO50106001
Tên HP: Phân tích thị trường bất động sản (ECO501060)
Giảng viên : Nguyễn Lưu Bảo Đoan (Email: doannlb@ueh.edu.vn), Trịnh Hoài Đức (Email: ducth@ueh.edu.vn)	3	HPTC.II.PTBDS.RE.2	Tư	2 (07g10)	A211	13/03/2024->13/03/2024)	A	59C Nguyễn Đình Chiểu, P.Võ Thị Sáu, Q.3, TP.HCM
Tư	2 (07g10)	NGHI	20/03/2024->27/03/2024)	Online	Học online
Bảy	2 (07g10)	A104a	06/04/2024->11/05/2024)	A	59C Nguyễn Đình Chiểu, P.Võ Thị Sáu, Q.3, TP.HCM
Bảy	2 (07g10)	NGHI	20/04/2024->20/04/2024)	Online	Học online
8	Mã LHP: 24D1MAR50303301
Tên HP: Marketing kỹ thuật số (MAR503033)
Giảng viên : Lê Thị Huệ Linh (Email: linhle.mar@ueh.edu.vn)	3	HPTC.II.PTBDS.RE.2	Tư	8 (12g45)	LMS	08/05/2024->08/05/2024)	LMS	Học trên hệ thống LMS
Tư	8 (12g45)	N2-301	13/03/2024->15/05/2024)	UEH Nguyễn Văn Linh - N2	Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM
9	Mã LHP: 24D1MAN50201702
Tên HP: Lập kế hoạch kinh doanh (MAN502017)
Giảng viên : Ngô Diễm Hoàng (Email: ngodiemhoang@ueh.edu.vn)	3	HPTC.I.BV.2	Năm	2 (07g10)	C(1.03)	14/03/2024->16/05/2024)	C	91 đường 3 tháng 2 Quận 10 TP.HCM
Năm	2 (07g10)	LMS	21/03/2024->21/03/2024)	LMS	Học trên hệ thống LMS
© Copyright 2020 PSC. ĐẠI HỌC UEH
Trụ sở: 59C Nguyễn Đình Chiểu, quận 3, TP. Hồ Chí Minh
Điện thoại: 84.28.38295299 - Fax: 84.28.38250359
Email: info@ueh.edu.vn`)
	exp := `{ Year: 2024, Semester: HKD, Courses: { { CourseName: Phí Bảo hiểm Y tế và Tai nạn, CourseCode: INS536001 }, { CourseName: Phát triển bất động sản nâng cao, CourseCode: ECO501179, Lecturer: Lê Nguyệt Trân, LecturerEmail: tranln@ueh.edu.vn, ClassCode: HPTC.I.PTBDS.RE.2, Credits: , Schedule: { { Day: Hai, Session: 8 (12g45), Room: N2-308, StartDate: 18/03/2024, EndDate: 3/05/2024, Campus: UEH Nguyễn Văn Linh - N2, Address: Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM}, { Day: Ba, Session: 2 (07g10), Room: N2-311, StartDate: 05/03/2024, EndDate: 5/03/2024, Campus: UEH Nguyễn Văn Linh - N2, Address: Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM} } }, { CourseName: Khởi nghiệp kinh doanh, CourseCode: BUS533002, Lecturer: Cao Văn Tiến, LecturerEmail: vantien@ueh.edu.vn, ClassCode: C.KNKD.A116_007, Credits: , Schedule: { { Day: Hai, Session: 13 (17g45), Room: A116, StartDate: 08/04/2024, EndDate: 5/04/2024, Campus: A, Address: 59C Nguyễn Đình Chiểu, P.Võ Thị Sáu, Q.3, TP.HCM}, { Day: Tư, Session: 13 (17g45), Room: A116, StartDate: 10/04/2024, EndDate: 0/04/2024, Campus: A, Address: 59C Nguyễn Đình Chiểu, P.Võ Thị Sáu, Q.3, TP.HCM}, { Day: Sáu, Session: 13 (17g45), Room: A116, StartDate: 12/04/2024, EndDate: 2/04/2024, Campus: A, Address: 59C Nguyễn Đình Chiểu, P.Võ Thị Sáu, Q.3, TP.HCM} } }, { CourseName: Sinh hoạt lớp hk6.1[7], CourseCode: ADM535024, Lecturer: Nguyễn Lưu Bảo Đoan, LecturerEmail: doannlb@ueh.edu.vn, ClassCode: RE002, Schedule: { { Day: Ba, Session: 2 (07g10), Room: ONLINE, StartDate: 14/05/2024, EndDate: 4/05/2024, Campus: Online, Address: Học online} } }, { CourseName: Hệ thống thông tin địa lý trong quy hoạch và quản lý đô thị, CourseCode: ARC512047, Lecturer: Nguyễn Minh Quang, LecturerEmail: quangnm@ueh.edu.vn), Phạm Nguyễn Hoà, ClassCode: HPTC.II.PTBDS.RE.2, Credits: , Schedule: { { Day: Ba, Session: 2 (07g10), Room: N1-403, StartDate: 23/01/2024, EndDate: 3/01/2024, Campus: UEH Nguyễn Văn Linh - N1, Address: Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM}, { Day: Tư, Session: 2 (07g10), Room: N2-210, StartDate: 17/01/2024, EndDate: 6/03/2024, Campus: UEH Nguyễn Văn Linh - N2, Address: Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM}, { Day: Tư, Session: 2 (07g10), Room: NGHI, StartDate: 03/01/2024, EndDate: 0/01/2024, Campus: Online, Address: Học online}, { Day: Tư, Session: 2 (07g10), Room: ONLINE, StartDate: 31/01/2024, EndDate: 1/02/2024, Campus: Online, Address: Học online} } }, { CourseName: Doanh nghiệp bất động sản, CourseCode: ECO501180, Lecturer: Nguyễn Thị Hồng Thu, LecturerEmail: thunguyen@ueh.edu.vn), Nguyễn Thị Hoàng Thù, ClassCode: HPTC.I.PTBDS.RE.2, Credits: , Schedule: { { Day: Ba, Session: 8 (12g45), Room: B2-202, StartDate: 12/03/2024, EndDate: 7/05/2024, Campus: B2, Address: 279 Nguyễn Tri Phương P.5 Q.10 TP.HCM (Khu B2)}, { Day: Ba, Session: 8 (12g45), Room: LMS, StartDate: 14/05/2024, EndDate: 4/05/2024, Campus: LMS, Address: Học trên hệ thống LMS} } }, { CourseName: Phân tích thị trường bất động sản, CourseCode: ECO501060, Lecturer: Nguyễn Lưu Bảo Đoan, LecturerEmail: doannlb@ueh.edu.vn), Trịnh Hoài Đứ, ClassCode: HPTC.II.PTBDS.RE.2, Credits: , Schedule: { { Day: Tư, Session: 2 (07g10), Room: A211, StartDate: 13/03/2024, EndDate: 3/03/2024, Campus: A, Address: 59C Nguyễn Đình Chiểu, P.Võ Thị Sáu, Q.3, TP.HCM}, { Day: Tư, Session: 2 (07g10), Room: NGHI, StartDate: 20/03/2024, EndDate: 7/03/2024, Campus: Online, Address: Học online}, { Day: Bảy, Session: 2 (07g10), Room: A104a, StartDate: 06/04/2024, EndDate: 1/05/2024, Campus: A, Address: 59C Nguyễn Đình Chiểu, P.Võ Thị Sáu, Q.3, TP.HCM}, { Day: Bảy, Session: 2 (07g10), Room: NGHI, StartDate: 20/04/2024, EndDate: 0/04/2024, Campus: Online, Address: Học online} } }, { CourseName: Marketing kỹ thuật số, CourseCode: MAR503033, Lecturer: Lê Thị Huệ Linh, LecturerEmail: linhle.mar@ueh.edu.vn, ClassCode: HPTC.II.PTBDS.RE.2, Credits: , Schedule: { { Day: Tư, Session: 8 (12g45), Room: LMS, StartDate: 08/05/2024, EndDate: 8/05/2024, Campus: LMS, Address: Học trên hệ thống LMS}, { Day: Tư, Session: 8 (12g45), Room: N2-301, StartDate: 13/03/2024, EndDate: 5/05/2024, Campus: UEH Nguyễn Văn Linh - N2, Address: Khu chức năng số 15, Đô thị mới Nam TP, Xã Phong Phú, Huyện Bình Chánh, TP.HCM} } }, { CourseName: Lập kế hoạch kinh doanh, CourseCode: MAN502017, Lecturer: Ngô Diễm Hoàng, LecturerEmail: ngodiemhoang@ueh.edu.vn, ClassCode: HPTC.I.BV.2, Credits: , Schedule: { { Day: Năm, Session: 2 (07g10), Room: C(1.03), StartDate: 14/03/2024, EndDate: 6/05/2024, Campus: C, Address: 91 đường 3 tháng 2 Quận 10 TP.HCM}, { Day: Năm, Session: 2 (07g10), Room: LMS, StartDate: 21/03/2024, EndDate: 1/03/2024, Campus: LMS, Address: Học trên hệ thống LMS} } } } }`

	if result.String() != exp {
		t.Fatalf("Expected %s but got %s", exp, result.String())
	}
}