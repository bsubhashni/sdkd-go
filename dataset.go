package main

import (
	"bytes"
	"fmt"
	"strconv"
)

type DatasetIterator interface {
	Start()
	Advance() bool
	Done() bool
	Key() string
	Value() string
}

type DatasetSeededIterator struct {
	spec   DS
	curk   string
	curv   string
	curidx int
}

func getDatasetIterator(spec DS) DatasetIterator {
	ds := new(DatasetSeededIterator)
	ds.SetSpec(spec)
	return ds
}

func (ds *DatasetSeededIterator) SetSpec(spec DS) {
	ds.spec = spec
}

func (ds *DatasetSeededIterator) Start() {
	ds.curidx = 0
	ds.initData()
}

func (ds *DatasetSeededIterator) fillRepeat(size int, seed string) string {
	filler := ds.spec.Repeat + strconv.Itoa(ds.curidx)
	base := seed + filler
	buf := new(bytes.Buffer)

	buf.Write([]byte(base))
	for buf.Len() < size {
		_, err := buf.Write([]byte(filler))
		if err != nil {
			fmt.Printf("Cannot write to dataset buffer: %v\n", err)
		}
	}
	return buf.String()
}

func (ds *DatasetSeededIterator) initData() {
	ds.curidx++
	ds.curk = ds.fillRepeat(ds.spec.KSize, ds.spec.KSeed)
	ds.curv = ds.fillRepeat(ds.spec.VSize, ds.spec.VSeed)
}

func (ds *DatasetSeededIterator) Advance() bool {
	if ds.spec.Continuous && ds.curidx > ds.spec.Count {
		ds.curidx = 0
	}
	ds.initData()
	return true
}

func (ds *DatasetSeededIterator) Done() bool {
	if ds.spec.Continuous == true {
		return false
	}

	if ds.curidx >= ds.spec.Count {
		return true
	}

	return false
}

func (ds *DatasetSeededIterator) Key() string {
	return ds.curk
}

func (ds *DatasetSeededIterator) Value() string {
	return ds.curv
}
