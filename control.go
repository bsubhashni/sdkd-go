package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type Control struct {
	Conn        net.Conn
	OutBuf      []byte
	InBuf       []byte
	GotRequest  chan bool
	ShouldFlush chan bool
	Quit        chan bool
}

func (controller *Control) ReadRequest() bool {
	rdr := bufio.NewReader(controller.Conn)

	for {
		buf := make([]byte, 1024)
		bytesRead, err := rdr.Read(buf)

		if err != nil {
			if err == io.EOF {
				return false
			} else {
				log.Fatalf("Error reading from Control Socket %v \n", err)
			}
		}

		if bytesRead == 0 {
			fmt.Printf("Remote has closed the connection \n")
			return false
		}
		fmt.Printf("Reading %d bytes from control socket \n", bytesRead)
		controller.GotRequest <- true
		controller.InBuf = buf[:bytesRead]
		return true
	}
}

func (controller *Control) ProcessRequest() {
	buf := controller.InBuf

	fmt.Printf("Got Message %s", string(buf))

	var req RequestCommand
	var res ResponseCommand
	if err := json.Unmarshal(buf, &req); err != nil {
		fmt.Printf("Cannot unmarshal command %v %v \n", err, req)
	}

	if req.Command == "INFO" {
		res.Command = req.Command
		res.ResData = InfoResponse{}
		b, err := json.Marshal(res)
		if err != nil {
			fmt.Printf("Unable to marshal info response \n")
		} else {
			controller.OutBuf = b
		}

	}

    if req.Command == "GOODBYE" {
    }

	controller.ShouldFlush <- true
}

func (controller *Control) WriteResponse() {
	buf := controller.OutBuf

	out := string(buf) + "\n"

	for {
		bytesWritten, err := controller.Conn.Write([]byte(out))

		if err != nil {
			log.Fatalf("Writing to control socket errored %v", err)
		}

		if bytesWritten == len([]byte(out)) {
			fmt.Printf("Successfully wrote %s \n", out)
			break
		}
	}
}

func (controller *Control) RequestHandler() {
	for {
		select {
		case <-controller.GotRequest:
			go controller.ProcessRequest()
		case <-controller.ShouldFlush:
			go controller.WriteResponse()
		default:
		}
	}
}

func (controller *Control) Start(conn net.Conn) {
	fmt.Println("Starting Controller \n")

	controller.Conn = conn
	controller.GotRequest = make(chan bool)
	controller.ShouldFlush = make(chan bool)

	go controller.ReadRequest()
	go controller.RequestHandler()

}
