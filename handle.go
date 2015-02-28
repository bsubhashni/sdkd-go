package main

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/go-couchbase"
	"github.com/couchbaselabs/gocb"
	"log"
	"net/url"
	"strconv"
)

type Handle interface {
	Init(DatasetIterator, *Options)
	CreateNewCouchbaseConnection(string, int, string, string, string) error
	DsMutate()
	DsGet()
	DsViewLoad()
	DsViewQuery()
	GetResult() *ResultResponse
	Cancel()
}

type Handle_v1 struct {
	couchbaseBucket *couchbase.Bucket
	DsIter          DatasetIterator
	rs              *ResultSet
	DoCancel        bool
	Schema          ViewSchema
}

type Handle_v2 struct {
	bucket   *gocb.Bucket
	DsIter   DatasetIterator
	rs       *ResultSet
	DoCancel bool
	Schema   ViewSchema
}

func (handle *Handle_v1) Init(ds DatasetIterator, opts *Options) {
	handle.DsIter = ds
	handle.rs = new(ResultSet)
	handle.rs.Initialize()
    handle.Schema = opts.VSchema
}


func (handle *Handle_v1) CreateNewCouchbaseConnection(hostname string, port int,
	bucket string, username string, password string) (err error) {

	var connStr string
	if password != "" {
		userinfo := url.UserPassword(username, password)
		u := &url.URL{
			Scheme: "http",
			User:   userinfo,
			Host:   hostname + ":" + strconv.Itoa(port),
		}
		connStr = u.String()
	} else {
		u := &url.URL{
			Scheme: "http",
			Host:   hostname + ":" + strconv.Itoa(port),
		}
		connStr = u.String()
	}

	c, err := couchbase.Connect(connStr)
	if err != nil {
		return err
	}

	p, err := c.GetPool(bucket)
	if err != nil {
		return err
	}

	if password != "" {
		handle.couchbaseBucket, err = p.GetBucket(bucket)
		if err != nil {
			return err
		}
	} else {
		handle.couchbaseBucket, err = p.GetBucketWithAuth(bucket, username, password)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Successfully instantiated connection \n")
	return nil
}

func (handle *Handle_v1) DsMutate() {
	dsIter := handle.DsIter
	handle.DoCancel = false

	for dsIter.Start(); dsIter.Done() == false && handle.DoCancel == false; dsIter.Advance() {
		key := dsIter.Key()
		val := dsIter.Value()

		handle.rs.MarkBegin()
		err := handle.couchbaseBucket.Set(key, 0, val)
		if err != nil {
			log.Fatalf("Cannot set items: %v key %v value %v \n", err, key, val)

		}
	}
}

func (handle *Handle_v1) DsGet() {
	dsIter := handle.DsIter
	handle.DoCancel = false

	for dsIter.Start(); dsIter.Done() == false && handle.DoCancel == false; dsIter.Advance() {
		key := dsIter.Key()
		var val string

		handle.rs.MarkBegin()
		err := handle.couchbaseBucket.Get(key, val)
		if err != nil {
			log.Fatalf("Cannot set items: %v key %v value %v \n", err, key, val)

		}
	}
}

func (handle *Handle_v1) DsViewQuery() {
	log.Fatalf("Not implemented in legacy couchbase sdk")
}

func (handle *Handle_v1) DsViewLoad() {
	log.Fatalf("Not implemented in legacy couchbase sdk")
}

func (handle *Handle_v1) Cancel() {
	handle.DoCancel = true
}

func (handle *Handle_v1) GetResult() *ResultResponse {
	res := new(ResultResponse)
	handle.rs.ResultsJson(res)
	return res
}

func (handle *Handle_v2) Init(dsIter DatasetIterator, opts *Options) {
	handle.DsIter = dsIter
	handle.rs = new(ResultSet)
	handle.rs.Initialize()
	handle.rs.Options = opts
}

func (handle *Handle_v2) CreateNewCouchbaseConnection(hostname string, port int,
	bucket string, username string, password string) (err error) {
	connStr := "couchbase://" + hostname

	c, err := gocb.Connect(connStr)
	if err != nil {
		return err
	}

	handle.bucket, err = c.OpenBucket(bucket, password)
	if err != nil {
		return err
	}
	return nil
}

func (handle *Handle_v2) DsMutate() {
	dsIter := handle.DsIter
	handle.DoCancel = false

	for dsIter.Start(); dsIter.Done() == false && handle.DoCancel == false; dsIter.Advance() {
		key := dsIter.Key()
		val := dsIter.Value()

		handle.rs.MarkBegin()

		_, err := handle.bucket.Upsert(key, val, 0)
		if err != nil {
			//log.Fatalf("Cannot set items using handle v2 %v %v %v\n", err, key, val)
			handle.rs.setResCode(1, key, val, "")
		} else {
			handle.rs.setResCode(0, key, val, "")
		}
	}
}

func (handle *Handle_v2) DsGet() {
	dsIter := handle.DsIter
	handle.DoCancel = false

	for dsIter.Start(); dsIter.Done() == false && handle.DoCancel == false; dsIter.Advance() {
		key := dsIter.Key()
		expectedVal := dsIter.Value()
		var v string

		handle.rs.MarkBegin()

		_, _, err := handle.bucket.Get(key, v)

		if err != nil {
			log.Fatalf("Cannot get items using handle v2 %v %v \n", err, key)
			handle.rs.setResCode(1, key, v, expectedVal)
		} else {
			handle.rs.setResCode(0, key, v, expectedVal)
		}
	}
}

func (handle *Handle_v2) DsViewLoad() {
	dsIter := handle.DsIter
	handle.DoCancel = false

	ii := 0
	for dsIter.Start(); dsIter.Done() == false && handle.DoCancel == false; dsIter.Advance() {
		key := dsIter.Key()
		handle.Schema.KIdent = dsIter.Key()
		handle.Schema.KSequence = ii

		handle.rs.MarkBegin()

		b, err := json.Marshal(handle.Schema)
		if err != nil {
			log.Fatalf("Unable to marshal schema for view load")
		}

		_, err = handle.bucket.Upsert(key, b, 0)

		if err != nil {
			log.Fatalf("Cannot get items using handle v2 %v %v \n", err, key)
			handle.rs.setResCode(1, key, "", "")
		} else {
			handle.rs.setResCode(0, key, "", "")
		}
		ii++
	}

}

func (handle *Handle_v2) DsViewQuery() {
}

func (handle *Handle_v2) GetResult() *ResultResponse {
	res := new(ResultResponse)
	handle.rs.ResultsJson(res)
	return res
}

func (handle *Handle_v2) Cancel() {
	handle.DoCancel = true
}
