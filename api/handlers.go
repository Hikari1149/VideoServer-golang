package main


import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
	"videoServer/api/dbops"
	"videoServer/api/defs"
	"videoServer/api/session"
)

func CreateUser (w http.ResponseWriter, r *http.Request, p httprouter.Params){
	res,_:=ioutil.ReadAll(r.Body)
	ubody :=&defs.UserCredential{
		UserName: "",
		Pwd:      "",
	}
	if err:= json.Unmarshal(res,ubody);err!=nil{
		sendErrorResponse(w,defs.ErrorRequestBodyParseFailed)
		return
	}
	if err:=dbops.AddUserCredential(ubody.UserName,ubody.Pwd); err!=nil{
		sendErrorResponse(w,defs.ErrorDBError)
		return
	}
	//create session
	id:=session.GenerateNewSessionId(ubody.UserName)
	su:=&defs.SignedUp{
		Success:   true,
		SessionId: id,
	}
	if resp,err:= json.Marshal(su);err!=nil{
		sendErrorResponse(w,defs.ErrorInternalFaults)
		return
	}else{
		sendNormalResponse(w,string(resp),201)
	}

}

func Login (w http.ResponseWriter, r *http.Request, p httprouter.Params){
	uName:=p.ByName("user_name")
	io.WriteString(w,uName)
}



func GetUserInfo(w http.Response,r *http.Request,p httprouter.Params){

}



