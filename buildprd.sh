#! /bin/bash

#build web and other services

cd ~/go/src/videoServer/api
env GOOS=linux GOARCH=amd64 go build -o ../bin/api

cd ~/go/src/videoServer/scheduler
env GOOS=linux GOARCH=amd64 go build -o ../bin/scheduler

cd ~/go/src/videoServer/streamServer
env GOOS=linux GOARCH=amd64 go build -o ../bin/streamServer

cd ~/go/src/videoServer/web
env GOOS=linux GOARCH=amd64 go build -o ../bin/web