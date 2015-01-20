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
	parent      *Sdkd
	Conn        net.Conn
	OutBuf      []byte
	InBuf       []byte
	GotRequest  chan bool
	ShouldFlush chan bool
	handle      Handle
	CloseConn   chan bool
	Quit        chan bool
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
		fmt.Printf("New handle\n")
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

		worker.parent.Mutex.Lock()
		worker.parent.HandleMap[req.Handle] = worker
		worker.parent.Mutex.Unlock()
	}

	if req.Command == "CANCEL" {
		fmt.Printf("Cancel Handle \n")
		res.ResData = EmptyObject{}
	}

	if req.Command == "CLOSEHANDLE" {
		fmt.Printf("Close Handle\n")
		res.ResData = EmptyObject{}
		res.Status = 0
		worker.parent.Mutex.Lock()
		delete(worker.parent.HandleMap, req.ReqID)
		worker.parent.Mutex.Unlock()
	}

	//Create Dataset Iterator
	handle.Init(getDatasetIterator(req.CmdData.DS), &req.CmdData.Options)

	if req.Command == "MC_DS_MUTATE_SET" {
		handle.DsMutate()
		res.ResData = handle.GetResult()
		fmt.Printf("%v", res.ResData)
	}

	b, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("Unable to marshal %s response", res.Command)
	}
	worker.OutBuf = b
	fmt.Printf("Worker out buffer %s \n", string(b))

	worker.ShouldFlush <- true

	if req.Command == "CLOSEHANDLE" {
		worker.Quit <- true
	}
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
		case <-worker.Quit:
			worker.Conn.Close()
			break
		default:
		}
	}
}

func (worker *Worker) Start(conn net.Conn) {
	fmt.Println("Starting new worker \n")
	var h Handle_v2

	worker.Conn = conn
	worker.handle = &h
	worker.GotRequest = make(chan bool)
	worker.ShouldFlush = make(chan bool)
	worker.Quit = make(chan bool)

	go worker.ReadRequest()
	go worker.RequestHandler()

}
