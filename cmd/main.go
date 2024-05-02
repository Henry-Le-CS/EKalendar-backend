package main

import (
	"e-calendar/cmd/modules/calendar"
	"e-calendar/cmd/modules/processor"
	"log"

	"github.com/Henry-Le-CS/gost"
)

func main() {
	modules := []gost.IModule{
		processor.ProcessorModule(),
		calendar.CalendarModule(),
	}

	s := gost.NewServer("localhost:8080", modules)

	if err := s.Start(); err != nil {
		log.Fatal(err)
		panic(err)
	}
}