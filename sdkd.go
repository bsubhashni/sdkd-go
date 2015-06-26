package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Sdkd struct {
	Port          int
	ShouldPersist bool
	HandleMap     map[int]*Worker
	Mutex         sync.Mutex
	Handle        int
	logger        *Logger
}

/* SDKD driver begins accept new connections */
func (sdkd *Sdkd) Start() (err error) {
	connCount := 0
	sdkd.HandleMap = make(map[int]*Worker)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", sdkd.Port))
	if err != nil {
		log.Fatalf("Cannot listen on port %v", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Cannot accept new connection %v", err)
		}

		if connCount == 0 {
			/** First connection is special - control socket  **/
			connCount++
			control := new(Control)
			control.parent = sdkd
			control.logger = sdkd.logger
			go control.Start(conn)
		} else {
			worker := new(Worker)
			worker.parent = sdkd
			worker.logger = sdkd.logger
			go worker.Start(conn)
		}
	}
}
