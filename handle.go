package main

import (
	"fmt"
	"github.com/couchbaselabs/go-couchbase"
	"log"
	"net/url"
)

type Handle struct {
	couchbaseBucket *couchbase.Bucket
	DsIter          DatasetIterator
}

func (handle *Handle) CreateNewCouchbaseConnection(hostname string, port int,
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

func (handle *Handle) dsMutate() {
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
