package main


import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"videoServer/api/session"
)

type middleWareHandler struct{
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler{
	m:=middleWareHandler{}
	m.r=r
	return m
}
func (m middleWareHandler) ServeHTTP (w http.ResponseWriter,r *http.Request){
	//check session
	validateUserSession(r)
	m.r.ServeHTTP(w,r)
}



func RegisterHandlers() *httprouter.Router{
	router:= httprouter.New()

	router.POST("/user",CreateUser)

	router.POST("/user/:user_name",Login)


	router.GET("/user/:username",GetUserInfo)

	return router
}

func Prepare(){
	session.LoadSessionFromDB()
}
func main(){
	Prepare()
	r:=RegisterHandlers()
	mh:=NewMiddleWareHandler(r)
	http.ListenAndServe(":8000",mh)
}
