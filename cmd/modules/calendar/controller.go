package calendar

import (
	"e-calendar/cmd/common"
	calender_services "e-calendar/cmd/modules/calendar/services"
	processor_srv "e-calendar/cmd/modules/processor/services"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Henry-Le-CS/gost"
	"golang.org/x/oauth2"
)

func CalendarController() *gost.Controller {
	router := gost.DeclareRouter()

	router.Add(gost.DeclareRouteHandler(
		"POST",
		"",
		createCalendar,
	))

	router.Add(gost.DeclareRouteHandler(
		"POST",
		"/add",
		AddGoogleCalendar,
	))

	controllers := gost.DeclareController(gost.ControllerArgs{
		Prefix: "/calendar",
		Router: router,
	})

	return controllers
}

func createCalendar(w http.ResponseWriter, r *http.Request) {
	var calendarRequestDto CalendarRequestDto
	err := json.NewDecoder(r.Body).Decode(&calendarRequestDto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isValid:= validateCalendarRequestDto(calendarRequestDto, w); !isValid {
		return
	}

	processorService := processor_srv.NewProcessorService("ueh")
	calendarService := calender_services.NewCalendarService(calendarRequestDto.University, processorService)

	if calendarService == nil {
		common.RaiseBadRequest(w, "University is not supported")
		return
	}

	text := strings.Join(calendarRequestDto.Texts, "\n")

	calendar, err := calendarService.CreateCalendar(text)

	if err != nil {
		common.RaiseBadRequest(w, err.Error())
		return
	}

	res := common.GenerateResponse(calendar, "")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func validateCalendarRequestDto(dto CalendarRequestDto, w http.ResponseWriter) bool {
	if dto.Semester == "" {
		common.RaiseBadRequest(w, "Semester is required")
		return false
	} else if dto.University == "" {
		common.RaiseBadRequest(w, "University is required")
		return false
	} else if dto.Texts == nil || len(dto.Texts) == 0 {
		common.RaiseBadRequest(w, "Text is required")
		return false
	}

	return true
}

// ============= Add new calendar to Google Calendar =============

type ReqToken struct {
	AccessToken string
	TokenType string
}
type AddGCalDto struct {
	Ics string
	Token ReqToken
	CalendarName string
}

func AddGoogleCalendar(w http.ResponseWriter, r *http.Request) {
	var body AddGCalDto

	err := json.NewDecoder(r.Body).Decode(&body)
	
	if err != nil {
		common.RaiseBadRequest(w, err.Error())
		return
	}
	
	gcalService := NewGoogleCalendarService()
	
	err = gcalService.UploadNewCalendar(body.Ics, body.CalendarName, &oauth2.Token{
		AccessToken: body.Token.AccessToken,
		TokenType: body.Token.TokenType,
	})

	if err != nil {
		common.RaiseBadRequest(w, err.Error())
		return
	}

	res := common.GenerateResponse("Calendar added successfully", "")
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
