#!/bin/bash
go get -t -d -v ./... && go build -v ./...

go build -o main.go page.go