package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandler ()*httprouter.Router{
	router:=httprouter.New()

	router.GET("/",homeHandler)
	router.POST("/",homeHandler)
	router.GET("/userhome",userHomeHandler)
	router.POST("/userhome",userHomeHandler)

	////转发请求
	router.POST("/api",apiHandler)

	//
	router.GET("/videos/:vid-id",proxyVideoHandler)

	//proxyHandler
	router.POST("/upload/:vid-id",proxyHandler)


	////filer server
	router.ServeFiles("/statics/*filepath",http.Dir("./template")) //
	return router
}


func main(){
	r:=RegisterHandler()
	http.ListenAndServe(":8080",r)
}