package main

import (
	"fmt"
)

func main() {
	fmt.Printf("Starting SDKD on port 8050")

	sdkd := Sdkd{
		Port: "8050"}
	sdkd.Start()
}
