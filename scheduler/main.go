package main

import (
	"net/http"
	"videoServer/scheduler/taskRunner"
	"videoServer/src/github.com/julienschmidt/httprouter"
)

func RegisterHandlers() *httprouter.Router{
	router:=httprouter.New()

	router.GET("/video-delete-record/:vid-id",vidDelRecHandler)
	return router
}

/*
	user->api service-> delete video
	api service -> scheduler -> write video deletion
	timer
	timer->runner -> read vdr -> exec -> delete from folder
*/
func main(){
	go taskRunner.Start()
	r:=RegisterHandlers()
	http.ListenAndServe(":9001",r)
}