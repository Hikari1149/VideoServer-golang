package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func testPageHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	t,_:=template.ParseFiles("./videos/upload.html")
	t.Execute(w,nil)
}

func StreamHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	vid:=p.ByName("vid-id")
	vl:=VIDEO_DIR+vid
	video,err:=os.Open(vl)
	if err!=nil{
		log.Printf("open error %v",err)
		sendErrorResponse(w,http.StatusInternalServerError,"Internal Error")
		return
	}

	w.Header().Set("Content-Type","video/mp4")
	//ServeContent replies to the request using the content in the provided ReadSeeker.
	http.ServeContent(w,r,"",time.Now(),video) //二进制流形式传输给客户端
	defer video.Close()



	//targetUrl:="http://aliyun.XXX/videos"+p.ByName("vid-id")
	//http.Redirect(w,r,targetUrl,301)
}
func uploadHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	r.Body = http.MaxBytesReader(w,r.Body,MAX_UPLOAD_SIZE)
	if err:=r.ParseMultipartForm(MAX_UPLOAD_SIZE);err!=nil{
		sendErrorResponse(w,http.StatusBadRequest,"File Size is too big")
		return
	}

	file,_,err:=r.FormFile("file") //<form name="file">
	if err!=nil{
		sendErrorResponse(w,http.StatusInternalServerError,"Internal Error")
		return
	}

	data,err:=ioutil.ReadAll(file) //
	if err!=nil{
		log.Printf("Read file error")
		return
	}

	fn:=p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR+fn,data,0666) //
	if err!=nil{
		log.Printf("Write file error: %v",err)
		sendErrorResponse(w,http.StatusInternalServerError,"internal error")
		return
	}

	//oss
	ossfn:="videos/"+fn
	path:="./videos"+fn
	bn:="avenssi-videos2"
	ret :=UploadToOss(ossfn,path,bn)
	if !ret{
		sendErrorResponse(w,http.StatusInternalServerError,"failed")
		return
	}
	os.Remove(path)
	//oss_end

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w,"Uploaded successful")
}
