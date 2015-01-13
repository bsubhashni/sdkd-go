package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type Worker struct {
	Conn        net.Conn
	OutBuf      []byte
	InBuf       []byte
	GotRequest  chan bool
	ShouldFlush chan bool
	handle      *Handle
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
		worker.InBuf = buf[:bytesRead]
	}

}

func (worker *Worker) ProcessRequest() {
	buf := worker.InBuf
	fmt.Printf("Got Message %s", string(buf))

	var req RequestCommand
	var res ResponseCommand

	if err := json.Unmarshal(buf, &req); err != nil {
		fmt.Printf("Cannot unmarshal command %v %v \n", err, req)
	}

	res.Command = req.Command
	res.Handle = req.Handle
	res.ReqID = req.ReqID

	handle := worker.handle

	if req.Command == "NEWHANDLE" {
		res.ResData = EmptyObject{}

		var cmdData CommandData
		cmdData = req.CmdData

		if err := handle.CreateNewCouchbaseConnection(cmdData.Hostname,
			cmdData.Port,
			cmdData.Bucket,
			cmdData.Options.Username,
			cmdData.Options.Password); err != nil {
			fmt.Printf("Error establishing couchbase connection %v \n", err)
			res.Status = 1
		} else {
			res.Status = 0
		}
	}

	b, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("Unable to marshal %s response", res.Command)
	}
	worker.OutBuf = b
	fmt.Printf("Worker out buffer %s \n", string(b))
	worker.ShouldFlush <- true
}

func (worker *Worker) WriteResponse() {
	buf := worker.OutBuf

	out := string(buf) + "\n"
	for {
		bytesWritten, err := worker.Conn.Write([]byte(out))

		if err != nil {
			log.Fatalf("writing to worker socket errored")
		}
		if bytesWritten == len([]byte(out)) {
			fmt.Printf("Successfully wrote on worker socket %s \n", string(buf))
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
	fmt.Println("Starting new worker \n")

	worker.Conn = conn
	worker.handle = new(Handle)
	worker.GotRequest = make(chan bool)
	worker.ShouldFlush = make(chan bool)

	go worker.ReadRequest()
	go worker.RequestHandler()

}
