package processor_srv

import (
	"e-calendar/cmd/modules/processor"
	"strings"
)

type IProcessorService interface {
	ProcessFullPage(input string) (processor.CourseListDto, error)
	ProcessCourse(input string) processor.CourseDto
}

func NewProcessorService(university string) (IProcessorService) {
	switch strings.ToLower(university) {
		case "ueh":
			return NewUehProcessorService()
		default:
			return nil
	}
}