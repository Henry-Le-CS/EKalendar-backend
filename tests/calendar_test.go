package tests

import (
	calender_services "e-calendar/cmd/modules/calendar/services"
	"e-calendar/cmd/modules/processor"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestCreateCalendar(t *testing.T) {
	folder := "./input/calendar/tc1"
	input, err := os.ReadFile(folder + "/input.txt")

	if err != nil {
		t.Error(err)
	}

	processorService := processor.NewProcessorService()
	calendarService := calender_services.NewCalendarService("ueh")

	courseListDto := processorService.ProcessFullPage(string(input))

	res, err := calendarService.CreateCalendar(courseListDto)

	if err != nil {
		t.Error(err)
	}

	output, err := os.ReadFile(folder + "/output.txt")

	exp :=string(output)

	lines := strings.Split(res, "\n")
	expLines := strings.Split(exp, "\n")

	if len(lines) != len(expLines) {
		t.Errorf("Length mismatch, Res: %d != Exp: %d", len(lines), len(expLines))
		return
	}

	fmt.Println(len(lines), len(expLines))

	for i := 0; i < len(lines); i++ {
		lhs := strings.TrimSpace(lines[i])
		rhs := strings.TrimSpace(expLines[i])

		if lhs != rhs {
			t.Errorf("Expected:\n%s\nGot:\n%s", expLines[i], lines[i])
			return
		}
	}
}

func TestCreateSingleEvent(t *testing.T) {
	folder := "./input/calendar/tc2"
	input, err := os.ReadFile(folder + "/input.txt")

	if err != nil {
		t.Error(err)
		return
	}

	processorService := processor.NewProcessorService()
	calendarService := calender_services.NewCalendarService("ueh")

	courseListDto := processorService.ProcessFullPage(string(input))
	res, err := calendarService.CreateCalendar(courseListDto)

	if err != nil {
		t.Error(err)
	}

	output, err := os.ReadFile(folder + "/output.txt")

	if err != nil {
		t.Error(err)
		return
	}

	exp := string(output)

	lines := strings.Split(res, "\n")
	expLines := strings.Split(exp, "\n")

	if len(lines) != len(expLines) {
		t.Errorf("Length mismatch, Res: %d != Exp: %d", len(lines), len(expLines))
		return
	}

	for i := 0; i < len(lines); i++ {
		lhs := strings.TrimSpace(lines[i])
		rhs := strings.TrimSpace(expLines[i])

		if lhs != rhs {
			t.Errorf("Expected:\n%s\nGot:\n%s", expLines[i], lines[i])
			return
		}
	}
}