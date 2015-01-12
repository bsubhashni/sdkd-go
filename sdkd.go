package main

import (
    "net"
    "log"
)

type Sdkd struct {
    Port string
}

/* SDKD driver begins accept new connections */
func (sdkd *Sdkd) Start() (err error) {
    control := Control {}
    worker := Worker {}

    connCount := 0
    ln, err := net.Listen("tcp", ":8050")
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
            connCount++;
            go control.Start(conn)
        } else {
            go worker.Start(conn)
        }
    }
}
