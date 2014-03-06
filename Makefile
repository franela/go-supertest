export GOPATH=$(shell pwd)

test:
	go get github.com/franela/goblin
	go get github.com/franela/goreq
	go test -v
