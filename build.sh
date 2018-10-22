#!/bin/bash
go get -t -d -v ./... && go build -v ./...

go build -o rubberband application.go page.go