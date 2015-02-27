package main

import (
	"github.com/couchbaselabs/gocb"
)

type Value struct {
	Id    string `json:"Id"`
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func GetQuery(dname, vname string, parameters ViewQueryParameters) *gocb.ViewQuery {
	viewquery := gocb.NewViewQuery(dname, vname)

	if parameters.Limit > 0 {
		viewquery = viewquery.Limit(parameters.Limit)
	}

	if parameters.Stale == false {
		viewquery = viewquery.Stale(gocb.Before)
	} else if parameter.Stale == true {
		viewquery = viewquery.Stale(gocb.None)
	}

	if parameters.UpdateAfter == true {
		viewquery = viewquery.Stale(gocb.After)
	}

	if parameters.Skip > 0 {
		viewquery = viewquery.Skip(parameters.Skip)
	}

	return viewquery
}

func processResults(viewresults *gocb.ViewResults) error {
	var val Value
	for {
		success := viewresults.Next(&val)
		if success == false {
			err := viewresults.Close()
			return err
		}
	}
}
