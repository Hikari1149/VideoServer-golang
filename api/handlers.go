package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"videoServer/api/dbops"
	"videoServer/api/defs"
	"videoServer/api/session"
	"videoServer/api/utils"
)


//req post -> read req body -> parse json -> dbops -> send Res
// req get -> get qs -> dbops -> to json -> send res
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
	res,_:=ioutil.ReadAll(r.Body)
	ubody:=&defs.UserCredential{
		UserName: "",
		Pwd:      "",
	}
	if err:=json.Unmarshal(res,ubody);res!=nil{
		log.Printf("login parse json err: %v",err)
		sendErrorResponse(w,defs.ErrorRequestBodyParseFailed)
		return
	}
	//validate the request body
	pwd,err:=dbops.GetUserCredential(ubody.UserName)
	if err!=nil || len(pwd)==0 || pwd!=ubody.Pwd{
		sendErrorResponse(w,defs.ErrorNotAuthUser)
		return
	}
	//generate sessionId (according to userName)
	id:=session.GenerateNewSessionId(ubody.UserName)
	si:=&defs.SignedIn{
		Success:   true,
		SessionId: id,
	}
	if resp,err:=json.Marshal(si); err!=nil{
		sendErrorResponse(w,defs.ErrorInternalFaults)
	}else{
		sendNormalResponse(w,string(resp),200)
	}
}



func GetUserInfo(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	if !ValidateUser(w,r){
		log.Printf("Unauthorized user\n")
		return
	}
	uname:=p.ByName("username")
	u,err:=dbops.GetUser(uname)
	if err!=nil{
		log.Printf("Error in GetUserInfo: %s",err)
		sendErrorResponse(w,defs.ErrorDBError)
		return
	}
	ui:=&defs.UserInfo{Id:u.Id}
	if resp,err:=json.Marshal(ui);err!=nil{
		sendErrorResponse(w,defs.ErrorInternalFaults)
	}else{
		sendNormalResponse(w,string(resp),200)
	}
}

//add video info
func AddNewVideo(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	if !ValidateUser(w,r){
		log.Printf("Unauthorized User\n")
		return
	}
	res,_:=ioutil.ReadAll(r.Body)
	nvBody:=&defs.NewVideo{}
	if err:=json.Unmarshal(res,nvBody);err!=nil{
		log.Printf("%s",err)
		sendErrorResponse(w,defs.ErrorRequestBodyParseFailed)
		return
	}
	vi,err:=dbops.AddNewVideo(nvBody.AuthorId,nvBody.Name)
	log.Printf("Author id: %d, name: %s\n",nvBody.AuthorId,nvBody.Name)
	if err!=nil{
		log.Printf("Error in AddNewVideo:%s",err)
		sendErrorResponse(w,defs.ErrorDBError)
		return
	}
	if resp,err:=json.Marshal(vi);err!=nil{
		sendErrorResponse(w,defs.ErrorInternalFaults)
	}else{
		sendNormalResponse(w,string(resp),201)
	}
}
func ListAllVideo(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	if !ValidateUser(w,r){
		return
	}
	uname:=p.ByName("username")
	vs,err:=dbops.ListVideoInfo(uname,0,utils.GetCurrentTimestampSec())
	if err!=nil{
		log.Printf("Error in ListAllvideos: %s",err)
		sendErrorResponse(w,defs.ErrorDBError)
	}

	vsi:=&defs.VideosInfo{Videos:vs}
	if resp,err:=json.Marshal(vsi);err!=nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	}else{
		sendNormalResponse(w,string(resp),200)
	}
}

func DeleteVideo(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	if !ValidateUser(w,r){
		return
	}
	vid:=p.ByName("vid-id")
	err:=dbops.DeleteVideoInfo(vid)
	if err!=nil{
		log.Printf("Erro in DeleteVideo: %s",err)
		sendErrorResponse(w,defs.ErrorDBError)
		return
	}
	go utils.SendDeleteVideoRequest(vid)
	sendNormalResponse(w,"",204)
}
func PostComment(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	if !ValidateUser(w,r){
		return
	}
	reqBody,_:=ioutil.ReadAll(r.Body)
	cbody:=&defs.NewComment{}
	if err:=json.Unmarshal(reqBody,cbody);err!=nil{
		log.Printf("%s",err)
		sendErrorResponse(w,defs.ErrorRequestBodyParseFailed)
		return
	}
	vid:=p.ByName("vid-id")
	if err:=dbops.AddNewComments(vid,cbody.AuthorId,cbody.Content);err!=nil{
		log.Printf("Error in post comment: %s",err)
		sendErrorResponse(w,defs.ErrorDBError)
	}else{
		sendNormalResponse(w,"ok",201)
	}
}
func ShowComments(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	if !ValidateUser(w,r){
		return
	}
	vid:=p.ByName("vid-id")
	cm,err:=dbops.ListComments(vid,0,utils.GetCurrentTimestampSec())
	if err!=nil{
		log.Printf("Error in Show Comments: %s",err)
		sendErrorResponse(w,defs.ErrorDBError)
		return
	}
	cms:=&defs.Comments{Comments:cm}
	if resp,err:=json.Marshal(cms);err!=nil{
		sendErrorResponse(w,defs.ErrorInternalFaults)
	}else{
		sendNormalResponse(w,string(resp),200)
	}
}



