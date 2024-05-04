package main

import (
	"e-calendar/cmd/modules/calendar"
	"e-calendar/cmd/modules/processor"
	"log"
	"os"

	"github.com/Henry-Le-CS/gost"
	"github.com/joho/godotenv"
)

func main() {
	modules := []gost.IModule{
		processor.ProcessorModule(),
		calendar.CalendarModule(),
	}
	
	if err:= godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := ":8080"

	if port := os.Getenv("PORT"); port != "" {
		PORT = ":" + port
	}

	s := gost.NewServer(PORT, modules, &gost.ServerOptions{
		CorsOptions: gost.CorsOptions{
			AllowedOrigins: []string{"*"},
			AllowCredentials: true,
		},
	})


	if err := s.Start(); err != nil {
		log.Fatal(err)
		panic(err)
	}
}