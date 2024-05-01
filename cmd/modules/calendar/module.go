package calendar

import "github.com/Henry-Le-CS/gost"

func CalendarModule() *gost.Module {
	controllers := []gost.IController{CalendarControllers()}

	return &gost.Module{
		Name: "calendar",
		Controllers: controllers,
	}
}