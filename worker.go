package main

import (
	"fmt"
	"net"
    "bufio"
    "io"
    "log"
)

type Worker struct {
	Conn   net.Conn
	OutBuf []byte
	InBuf  []byte
    GotRequest chan bool
    ShouldFlush chan bool
}

func (worker *Worker) ReadRequest() {
	rdr := bufio.NewReader(worker.Conn)

	for {
		buf := make([]byte, 1024)
		bytesRead, err := rdr.Read(buf)

		if err != nil {
			if err == io.EOF {
				return
			} else {
				log.Fatalf("Error reading from Worker Socket %v", err)
			}
		}
		if bytesRead == 0 {
			log.Fatalf("Remote has closed the connection")
		}
		fmt.Printf("Reading %d bytes from worker socket \n", bytesRead)
		worker.GotRequest <- true
		worker.InBuf = buf
	}

}

func (worker *Worker) ProcessRequest() {
	buf := worker.InBuf
	fmt.Printf("Got Message %s", string(buf))
	worker.OutBuf = []byte("")
	worker.ShouldFlush <- true
}

func (worker *Worker) WriteResponse() {
	buf := worker.OutBuf
	for {
		bytesWritten, err := worker.Conn.Write(buf)

		if err != nil {
			log.Fatalf("writing to worker socket errored")
		}
		if bytesWritten == len(buf) {
			fmt.Printf("Successfully wrote %s", string(buf))
			break
		}
	}
}

func (worker *Worker) RequestHandler() {
	for {
		select {
		case <-worker.GotRequest:
			go worker.ProcessRequest()
		case <-worker.ShouldFlush:
			go worker.WriteResponse()
		default:
		}
	}
}

func (worker *Worker) Start(conn net.Conn) {
	fmt.Println("Starting new worker")

	worker.Conn = conn
	worker.GotRequest = make(chan bool)
	worker.ShouldFlush = make(chan bool)

	go worker.ReadRequest()
	go worker.RequestHandler()

}
