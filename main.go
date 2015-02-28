package main

import (
	"fmt"
    "flag"
)

func main() {
	port := flag.Int("Port", 8050, "Port for the SDKD to listen on")
	persist := flag.Bool("Persist", false, "Persist the SDKD[Do not kill on GOODBYE]")

	flag.Parse()

	fmt.Printf("Starting SDKD on port 8050")

	sdkd := Sdkd {
		Port:    *port,
		ShouldPersist: *persist }

	sdkd.Start()
}
