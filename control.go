package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
)

type Control struct {
	parent      *Sdkd
	Conn        net.Conn
	OutBuf      []byte
	InBuf       []byte
	GotRequest  chan bool
	ShouldFlush chan bool
	Quit        chan bool
	logger      *Logger
}

func (c *Control) ReadRequest() {
	rdr := bufio.NewReader(c.Conn)

	for {
		buf := make([]byte, 1024)
		bytesRead, err := rdr.Read(buf)

		if err != nil {
			if err == io.EOF {
				return
			} else {
				log.Fatalf("Error reading from Control Socket %v", err)
			}
		}

		if bytesRead == 0 {
			c.logger.Error(prettify() + "Remote has closed the connection")
			c.Quit <- true
		}
		c.logger.Debug(prettify()+"Reading %d bytes from control socket", bytesRead)
		c.GotRequest <- true
		c.InBuf = buf[:bytesRead]
	}
}

func (c *Control) ProcessRequest() {
	buf := c.InBuf

	var req RequestCommand
	var res ResponseCommand
	if err := json.Unmarshal(buf, &req); err != nil {
		c.logger.Error(prettify()+"Cannot unmarshal command %v %v", err, req)
	}

	c.logger.Debug(prettify()+"Got message %v", req.Command)
	res.Command = req.Command
	res.ReqID = req.ReqID

	if req.Command == "INFO" {
		var info InfoResponse
		//info.TIME = uint64(time.Now().Unix())
		res.ResData = info

	}

	if req.Command == "CANCEL" {
		//cancels the handle sent on request
		res.Handle = req.Handle
		res.ResData = EmptyObject{}
		c.parent.Mutex.Lock()
		if c.parent.HandleMap == nil {
			log.Fatalf("Cannot find the requested handle to cancel\n")
		} else {
			c.parent.HandleMap[req.Handle].handle.Cancel()
		}
		c.parent.Mutex.Unlock()
	}

	if req.Command == "GOODBYE" {
		//close all handles
		c.parent.Mutex.Lock()
		for handleid, worker := range c.parent.HandleMap {
			c.logger.Debug(prettify()+"Sending kill signal to handle worker %d", handleid)
			worker.Quit <- true
		}
		c.parent.Mutex.Unlock()
		res.ResData = EmptyObject{}

		if c.parent.ShouldPersist == false {
			os.Exit(0)
		}
	}

	b, err := json.Marshal(res)
	if err != nil {
		c.logger.Error(prettify() + "Unable to marshal info response")
	} else {
		c.OutBuf = b
	}
	c.ShouldFlush <- true
}

func (c *Control) WriteResponse() {
	buf := c.OutBuf

	out := string(buf) + "\n"

	for {
		bytesWritten, err := c.Conn.Write([]byte(out))

		if err != nil {
			log.Fatalf("Writing to control socket errored %v", err)
		}

		if bytesWritten == len([]byte(out)) {
			c.logger.Debug(prettify()+"Successfully wrote %s", out)
			c.OutBuf = []byte{}
			break
		}
	}
}

func (c *Control) RequestHandler() {
	for {
		select {
		case <-c.GotRequest:
			go c.ProcessRequest()
		case <-c.ShouldFlush:
			go c.WriteResponse()
		case <-c.Quit:
			break
		default:
		}
	}
}

func (c *Control) Start(conn net.Conn) {
	c.logger.Info("Starting Controller")

	c.Conn = conn
	c.GotRequest = make(chan bool)
	c.ShouldFlush = make(chan bool)

	go c.ReadRequest()
	go c.RequestHandler()

}
