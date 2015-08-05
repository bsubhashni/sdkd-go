package main

import (
	"flag"
	"fmt"
)

func main() {
	port := flag.Int("Port", 8050, "Port for the SDKD to listen on")
	persist := flag.Bool("Persist", false, "Persist the SDKD[Do not kill on GOODBYE]")
	handleType := flag.Int("Handle", 3, "Type of the sdk handle to use: "+
		"1. Legacy SDK "+
		"2. Synchronous "+
		"3. Asynchronous")
	logFile := flag.String("LogFile", "", "Log file for sdkd")
	logLevel := flag.Int("LogLevel", 2, "Log level for sdkd")

	flag.Parse()

	fmt.Printf("Starting SDKD on port 8050")

	logger := new(Logger)
	logger.Init(*logFile, *logLevel)

	sdkd := Sdkd{
		Port:          *port,
		ShouldPersist: *persist,
		Handle:        *handleType,
		logger:        logger}

	sdkd.Start()
}
