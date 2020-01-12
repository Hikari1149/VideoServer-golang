package session


import (
	"time"
	"sync"
	"videoServer/api/dbops"
	"videoServer/api/defs"
	"videoServer/api/utils"
)



var sessionMap *sync.Map

func init(){
	sessionMap = &sync.Map{}
}
func nowInMilli() int64{
	return time.Now().UnixNano()/100000
}


// retrieve all session from db -> load to map (normally redis)
func LoadSessionFromDB(){
	r,err:=dbops.RetrieveAllSessions()
	if err!=nil{
		return
	}
	r.Range(func(k,v interface{}) bool{
		ss:=v.(*defs.SimpleSession)
		sessionMap.Store(k,ss)
		return true
	})
}
func GenerateNewSessionId(un string) string{
	id,_:=utils.NewUUID()
	ct:=nowInMilli()
	ttl:=ct+ 30*60*1000  //ServerSide session valid time 30min

	ss:=&defs.SimpleSession{
		Username: un,
		TTL:      ttl,
	}
	sessionMap.Store(id,ss) //store in cache
	dbops.InsertSession(id,ttl,un) // store in db
	return id
}
func deleteExpiredSession(sid string){
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

func IsSessionExpired(sid string)(string,bool){
	ss,ok:=sessionMap.Load(sid)
	ct:=nowInMilli()
	if ok{
		if ss.(*defs.SimpleSession).TTL <  ct{
			deleteExpiredSession(sid)  //delete expired session
			return "",true
		}
		return ss.(*defs.SimpleSession).Username,false
	}else{
		ss,err:=dbops.RetrieveSession(sid)
		if err!=nil || ss==nil{
			return "",true
		}
		if ss.TTL<ct{
			deleteExpiredSession(sid)
			return "",true
		}
		sessionMap.Store(sid,ss)
		return ss.Username,false
	}
	return "", true
}


