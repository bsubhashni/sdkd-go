package main

import (
	"bufio"
	"encoding/json"
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
	logger      *Logger
}

func (w *Worker) ReadRequest() {
	rdr := bufio.NewReader(w.Conn)

	for {
		buf := make([]byte, 1024)
		bytesRead, err := rdr.Read(buf)

		if err != nil {
			if err == io.EOF {
				return
			} else {
				log.Fatalf(prettify()+"Error reading from Worker Socket %v", err)
			}
		}
		if bytesRead == 0 {
			log.Fatalf(prettify() + "Remote has closed the connection")
		}
		w.logger.Info("Reading %d bytes from worker socket", bytesRead)

		w.GotRequest <- true
		w.InBuf = buf[:bytesRead]
	}

}

func (w *Worker) ProcessRequest() {
	buf := w.InBuf
	w.logger.Info("Got Message on worker %s", string(buf))

	var req RequestCommand
	var res ResponseCommand

	if err := json.Unmarshal(buf, &req); err != nil {
		w.logger.Error(prettify()+"Cannot unmarshal command %v %v", err, buf)
	}

	res.Command = req.Command
	res.Handle = req.Handle
	res.ReqID = req.ReqID

	handle := w.handle

	if req.Command == NEWHANDLE {
		w.logger.Info("Creating a new handle")

		res.ResData = EmptyObject{}
		var cmdData CommandData
		cmdData = req.CmdData
		if err := handle.CreateNewCouchbaseConnection(cmdData.Hostname,
			cmdData.Port,
			cmdData.Bucket,
			cmdData.Options.Username,
			cmdData.Options.Password); err != nil {
			w.logger.Error(prettify()+"Error establishing couchbase connection %v", err)
			res.Status = 1
		} else {
			res.Status = 0
		}

		w.parent.Mutex.Lock()
		w.parent.HandleMap[req.Handle] = w
		w.parent.Mutex.Unlock()
	}

	if req.Command == CANCEL {
		w.logger.Info("Cancelling handle")
		res.ResData = EmptyObject{}
	}

	if req.Command == CLOSEHANDLE {
		w.logger.Info("Closing handle")
		res.ResData = EmptyObject{}
		res.Status = 0
		w.parent.Mutex.Lock()
		delete(w.parent.HandleMap, req.ReqID)
		w.parent.Mutex.Unlock()
	}

	//Create Dataset Iterator
	handle.Init(getDatasetIterator(req.CmdData.DS),
		&req.CmdData.Options,
		req.CmdData.VSchema,
		w.parent.logger)

	if req.Command == MC_DS_MUTATE_SET {
		handle.DsMutate()
		res.ResData = handle.GetResult()
	}

	if req.Command == MC_DS_GET {
		handle.DsGet()
		res.ResData = handle.GetResult()
	}

	if req.Command == CB_VIEW_LOAD {
		handle.DsViewLoad()
		res.ResData = handle.GetResult()
	}

	if req.Command == CB_VIEW_QUERY {
		handle.DsViewQuery(req.CmdData.DesignName, req.CmdData.ViewName, req.CmdData.ViewQueryParameters)
		res.ResData = handle.GetResult()
	}

	if req.Command == CB_N1QL_CREATE_INDEX {
		handle.DsN1QLCreateIndex()
		res.ResData = handle.GetResult()
	}

	if req.Command == CB_N1QL_QUERY {
		handle.DsN1QLQuery()
		res.ResData = handle.GetResult()
	}

	if res.ResData == nil {
		res.ResData = EmptyObject{}
	}

	b, err := json.Marshal(res)
	if err != nil {
		w.logger.Error(prettify()+"Unable to marshal %s response", res.Command)
	}
	w.OutBuf = b
	w.logger.Debug(prettify()+"Worker out buffer %s", string(b))

	w.ShouldFlush <- true

	if req.Command == "CLOSEHANDLE" {
		w.Quit <- true
	}
}

func (w *Worker) WriteResponse() {
	buf := w.OutBuf

	out := string(buf) + "\n"
	for {
		bytesWritten, err := w.Conn.Write([]byte(out))

		if err != nil {
			//log.Fatalf("writing to worker socket errored")
			return
		}
		if bytesWritten == len([]byte(out)) {
			w.logger.Debug(prettify()+"Successfully wrote on worker socket %s", string(buf))
			break
		}
	}
}

func (w *Worker) RequestHandler() {
	for {
		select {
		case <-w.GotRequest:
			go w.ProcessRequest()
		case <-w.ShouldFlush:
			go w.WriteResponse()
		case <-w.Quit:
			w.Conn.Close()
			break
		default:
		}
	}
}

func (w *Worker) Start(conn net.Conn) {
	w.logger.Debug("Starting new worker")

	if w.parent.Handle == 1 {
		var h Handle_v1
		w.handle = &h
	} else if w.parent.Handle == 2 {
		var h Handle_v2
		w.handle = &h
	} else {
		var h Handle_v3
		w.handle = &h
	}

	w.Conn = conn
	w.GotRequest = make(chan bool)
	w.ShouldFlush = make(chan bool)
	w.Quit = make(chan bool)

	go w.ReadRequest()
	go w.RequestHandler()

}
