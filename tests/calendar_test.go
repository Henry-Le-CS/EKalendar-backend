package tests

import (
	"e-calendar/cmd/modules/calendar"
	calender_services "e-calendar/cmd/modules/calendar/services"
	processor_srv "e-calendar/cmd/modules/processor/services"
	"fmt"
	"os"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

func TestCreateCalendar(t *testing.T) {
	folder := "./data/calendar/tc1"
	input, err := os.ReadFile(folder + "/input.txt")

	if err != nil {
		t.Error(err)
	}

	processorService := processor_srv.NewProcessorService("ueh")
	calendarService := calender_services.NewCalendarService("ueh", processorService)

	res, err := calendarService.CreateCalendar(string(input))

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
	folder := "./data/calendar/tc2"
	input, err := os.ReadFile(folder + "/input.txt")

	if err != nil {
		t.Error(err)
		return
	}

	processorService := processor_srv.NewProcessorService("ueh")
	calendarService := calender_services.NewCalendarService("ueh", processorService)

	res, err := calendarService.CreateCalendar(string(input))

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

func TestInsertCalendar(t *testing.T) {
	folder := "./data/calendar/tc1"

	input, err := os.ReadFile(folder + "/output.txt")

	if err != nil {
		t.Error(err)
		return
	}

	gcalService := calendar.NewGoogleCalendarService()

	// Make sure to change the tokens
	err = gcalService.UploadNewCalendar(string(input), "Test UEH 223", &oauth2.Token{
		AccessToken: "ya29.a0AXooCguwqLFeKwlu7xou6Klg0Sq49yeYDXFVxKaC2jUW2gIitkiraCG8yNTX_DTmiBsJZLfl7Q9W4dIYfQElNNpAT6xcLZGydHBLT-4WqhWG3xop_xpzXu9GtNWKMZg-piqF0xIbJokgXGh9Rs-2EVs3zC3cyDuK8fIaCgYKAQMSARESFQHGX2MiBAVFzYvaraUvnIEP5llyhA0170",
		TokenType: "Bearer",
	})

	if err != nil {
		t.Error(err)
	}
}