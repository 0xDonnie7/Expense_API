package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (app *application) server() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
