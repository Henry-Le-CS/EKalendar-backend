package processor

import "github.com/Henry-Le-CS/gost"

func ProcessorModule() *gost.Module {
	controllers := []gost.IController{ProcessorControllers()}

	module := gost.DeclareModule(gost.RegisterModuleDto{
		Name: "Processor module",
		Controllers: controllers,
	})

	return module
}