package main

import (
	"log"
	"log/syslog"
)

func main() {
	logger, err := syslog.New(syslog.LOG_SYSLOG, "main.go")

	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logger)
	log.Print("Reporting for duty!")
}
