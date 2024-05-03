package processor_srv

import (
	"e-calendar/cmd/modules/processor"
	"fmt"
	"strings"
)

type IProcessorService interface {
	ProcessFullPage(input string) processor.CourseListDto
	ProcessCourse(input string) processor.CourseDto
}

func NewProcessorService(university string) (IProcessorService, error) {
	switch strings.ToLower(university) {
		case "ueh":
			return NewUehProcessorService(), nil
		default:
			return nil, fmt.Errorf("university %s is not supported", university)
	}
}