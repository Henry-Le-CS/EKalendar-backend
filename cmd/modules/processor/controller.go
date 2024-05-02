package processor

import (
	"e-calendar/cmd/common"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Henry-Le-CS/gost"
)

func ProcessorControllers() *gost.Controller {
	router := gost.DeclareRouter()

	router.Add(gost.DeclareRouteHandler(
		"POST",
		"/",
		processCalendarFromText,
	))

	controllers := gost.DeclareController(gost.ControllerArgs{
		Prefix: "/processing",
		Router: router,
	})

	return controllers
}

type TestStruct struct {
	Abc string
}

func processCalendarFromText(w http.ResponseWriter, r *http.Request) {
	var calendarRequestDto CalendarRequestDto
	err := json.NewDecoder(r.Body).Decode(&calendarRequestDto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validateCalendarRequestDto(calendarRequestDto, w)

	fmt.Println(calendarRequestDto)
	res := common.GenerateResponse(TestStruct{Abc: "Hello"}, "")
	processorService := NewProcessorService()

	result := processorService.ProcessCourse(calendarRequestDto.Text)

	fmt.Println(result)
	w.Write(res)
}

func validateCalendarRequestDto(dto CalendarRequestDto, w http.ResponseWriter) error {
	if dto.Text == "" {
		common.RaiseBadRequest(w, "Text is required")
	} else if dto.Semester == "" {
		common.RaiseBadRequest(w, "Semester is required")
	}

	return nil
}