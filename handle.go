package main

import (
	"fmt"
	"github.com/brett19/gocouchbase"
	"github.com/couchbaselabs/go-couchbase"
	"log"
	"net/url"
)

type Handle interface {
	SetDatasetIterator(DatasetIterator)
	CreateNewCouchbaseConnection(string, int, string, string, string) error
	DsMutate()
}

type Handle_v1 struct {
	couchbaseBucket *couchbase.Bucket
	DsIter          DatasetIterator
}

type Handle_v2 struct {
	bucket *gocouchbase.Bucket
	DsIter DatasetIterator
}

func (handle *Handle_v1) SetDatasetIterator(ds DatasetIterator) {
	handle.DsIter = ds
}

func (handle *Handle_v1) CreateNewCouchbaseConnection(hostname string, port int,
	bucket string, username string, password string) (err error) {

	var connStr string
	if password != "" {
		userinfo := url.UserPassword(username, password)
		u := &url.URL{
			Scheme: "http",
			User:   userinfo,
			Host:   hostname + ":" + "8091",
		}
		connStr = u.String()
	} else {
		u := &url.URL{
			Scheme: "http",
			Host:   hostname + ":" + "8091",
		}
		connStr = u.String()
	}

	c, err := couchbase.Connect(connStr)
	if err != nil {
		return err
	}

	p, err := c.GetPool("default")
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

	for dsIter.Start(); dsIter.Advance(); dsIter.Done() {
		key := dsIter.Key()
		val := dsIter.Value()
		err := handle.couchbaseBucket.Set(key, 0, val)
		if err != nil {
			log.Fatalf("Cannot set items: %v key %v value %v \n", err, key, val)

		}
	}
}

func (handle *Handle_v1) DsGet() {
	dsIter := handle.DsIter

	for dsIter.Start(); dsIter.Advance(); dsIter.Done() {
		key := dsIter.Key()
		var val string
		err := handle.couchbaseBucket.Get(key, val)
		if err != nil {
			log.Fatalf("Cannot set items: %v key %v value %v \n", err, key, val)

		}
	}
}

func (handle *Handle_v2) SetDatasetIterator(dsIter DatasetIterator) {
    fmt.Printf("Setting data set iterator")
    if dsIter == nil {
        fmt.Printf("dataset iterator is nil")
    }
	handle.DsIter = dsIter
}

func (handle *Handle_v2) CreateNewCouchbaseConnection(hostname string, port int,
	bucket string, username string, password string) (err error) {
	connStr := "couchbase://" + hostname

	c, err := gocouchbase.Connect(connStr)
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


	for dsIter.Start(); dsIter.Advance(); dsIter.Done() {
		key := dsIter.Key()
		val := dsIter.Value()
		_, err := handle.bucket.Upsert(key, val, 0)
		if err != nil {
			log.Fatalf("Cannot set items using handle v2 %v %v %v\n", err, key, val)
		}
	}
}

func (handle *Handle_v2) DsGet() {
	dsIter := handle.DsIter

	for dsIter.Start(); dsIter.Advance(); dsIter.Done() {
		key := dsIter.Key()
		var v string
		_, _, err := handle.bucket.Get(key, v)

		if err != nil {
			log.Fatalf("Cannot get items using handle v2 %v %v \n", err, key)
		}
	}
}
