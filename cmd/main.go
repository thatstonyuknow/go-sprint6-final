package main

import (
    "log"
    "os"

    "github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {

    logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)


	srv := server.MyServer(logger)

	// Running server with recevied parameters
    logger.Println("Running server at 8080...")
    if err := srv.HTTPServer.ListenAndServe(); err != nil {
        logger.Fatal("Error while server start: ", err)
    }
}
