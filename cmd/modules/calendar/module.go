package calendar

import "github.com/Henry-Le-CS/gost"

func CalendarModule() *gost.Module {
	controllers := []gost.IController{CalendarController()}
	
	module := gost.DeclareModule(gost.RegisterModuleDto{
		Name:        "calendar",
		Controllers: controllers,
	})

	return module
}