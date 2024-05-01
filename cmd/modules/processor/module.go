package processor

import "github.com/Henry-Le-CS/gost"

func ProcessorModule() *gost.Module {
	controllers := []gost.IController{CalendarControllers()}

	module := gost.DeclareModule(gost.RegisterModuleDto{
		Name: "Calendar module",
		Controllers: controllers,
	})

	return module
}