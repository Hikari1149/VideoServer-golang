package main

import (
	"bytes"
	"encoding/json"
	"hikari/config"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var httpClient *http.Client

func init(){
	httpClient = &http.Client{}
}

func request(b *ApiBody,w http.ResponseWriter,r *http.Request){

	var resp *http.Response
	var err error

	u,_:=url.Parse(b.Url)
	u.Host=config.GetLbAddr()+":"+u.Port()
	//newUrl:=u.String()

	switch b.Method {
		case http.MethodGet:
				req,_:=http.NewRequest("GET",b.Url,nil)
				req.Header = r.Header
				resp,err = httpClient.Do(req) //
				if err!=nil {
					log.Printf("req error: %v", err)
					return
				}
				normalResponse(w,resp)
		case http.MethodPost:
			req,_:=http.NewRequest("POST",b.Url,bytes.NewBuffer([]byte(b.ReqBody )))
			req.Header=r.Header
			resp,err = httpClient.Do(req)
			if err!=nil{
				log.Printf("req post error: %v",err)
				return
			}
			normalResponse(w,resp)
		case http.MethodDelete:
			req,_:=http.NewRequest("DELETE",b.Url,nil)
			req.Header=r.Header
			resp,err = httpClient.Do(req)
			if err!=nil{
				log.Printf("Del req error: %v",err)
				return
			}
			normalResponse(w,resp)
		default:
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w,"Bad api request")
	}
}

func normalResponse(w http.ResponseWriter,r *http.Response){
	res,err:=ioutil.ReadAll(r.Body)
	if err!=nil{
		re,_:=json.Marshal(ErrorInternalFaults)
		w.WriteHeader(500)
		io.WriteString(w,string(re))
		return
	}
	w.WriteHeader(r.StatusCode)
	io.WriteString(w,string(res))
}