package main

import (
	"io"
	"log"
	"os"
)

func main() {
	// Logger
	InfoLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger := log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLogger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	logfile, err := os.Create("app.log")
	err = nil
	if err != nil {
		log.Fatal("Couldn't create logfile")
	} else {
		// Logging to both stdout and file
		multi := io.MultiWriter(logfile, os.Stdout)
		InfoLogger.SetOutput(multi)
		WarningLogger.SetOutput(multi)
		ErrorLogger.SetOutput(multi)
	}

	InfoLogger.Println("Starting TropesToGo...")

	// Setup config

	// Run services
	// app.Run(config)
}
