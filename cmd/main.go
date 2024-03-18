package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sivistrukov/vk-assigment/internal/app"
)

//	@title			VK assignment
//	@version		1.0
//	@description	The backend of the Filmotek application, which provides a REST API for managing the movie database.
//	@contact.name	Vladislav Strukov
//	@contact.email	sivistrukov@gmail.com
//	@host			localhost:8080
//	@BasePath		/api
//	@securityDefinitions.basic BasicAuth
func main() {
	log.Println("server starting...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("cant load .env file: %s\n", err.Error())
	}

	// TODO: remove
	err = godotenv.Load("../.env")
	if err != nil {
		log.Printf("cant load .env file: %s\n", err.Error())
	}

	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Println("server running")
		if err := app.Run(); err != nil {
			if err == http.ErrServerClosed {
				log.Println("server closed under request")
			} else {
				log.Fatalf("server closed unexpected: %s", err.Error())
			}
		}
	}()

	<-ctx.Done()

	stop()

	err = app.Shutdown()
	if err != nil {
		log.Println(err)
	}
}
