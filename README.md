#Go SDKD Implementation

This is a go based sdkd implementation, it is currently designed to work with legacy & new implementations of gocouchbase

##Dependencies:
* github.com/couchbaselabs/gocb
* github.com/couchbase/go-couchbase


##Prerequistes:
Go

##Build Steps:
make


This should create the executable sdkd-go. 

By default, the sdkd starts listening on port 8050.

##Options
* --Port Port for the sdkd to listen on
* --Persist Do not exit sdkd on goodbye
* --Handle 
  *  1 - Legacy go-couchbase 
  *  2 - Gocb sync ops 
  *  3 - Gocbcore async ops

