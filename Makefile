all:
	@go get github.com/couchbaselabs/go-couchbase
	@export GOPATH=$(PWD)
	@go build
