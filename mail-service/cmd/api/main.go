package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Mailer Mail
}

const webPort = "80"

func main() {

	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting mail service on port: ", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}

}

func createMail() Mail {

	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        os.Getenv("MAIL_PORT"),
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
	}

	return m
}
