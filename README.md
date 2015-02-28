#Go SDKD Implementation

This is a go based sdkd implementation, it is currently designed to work with legacy & new implementations of gocouchbase

##Dependencies:
github.com/couchbaselabs/gocb
github.com/couchbase/go-couchbase


##Prerequistes:
Go

##Build Steps:
make


This should create the executable sdkd-go. 

By default, the sdkd starts listening on port 8050. To listen on a different port 
./sdkd-go --Port 8050

