package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
)

type HomePage struct{
	Name string
}



func homeHandler(w http.ResponseWriter,r *http.Request,ps httprouter.Params){
	p:=&HomePage{Name:"avenssi"}
	t,e:=template.ParseFiles("./template/home.html")//html->parsed
	if e!=nil{
		log.Printf("Parsing template home.html error:%s",e)
		return
	}
	t.Execute(w,p) //// Execute applies a parsed template to the specified data object,
	// writing the output to wr.
	return
}