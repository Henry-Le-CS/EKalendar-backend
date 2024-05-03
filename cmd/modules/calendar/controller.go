package calendar

import (
	"e-calendar/cmd/common"
	calender_services "e-calendar/cmd/modules/calendar/services"
	processor_srv "e-calendar/cmd/modules/processor/services"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Henry-Le-CS/gost"
)

func CalendarController() *gost.Controller {
	router := gost.DeclareRouter()

	router.Add(gost.DeclareRouteHandler(
		"POST",
		"/",
		createCalendar,
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
