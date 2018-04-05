package main

import(
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"bytes"
	"io/ioutil"
	"github.com/gorilla/mux"
)

type Request struct {
	Port 		string `json:"port"`
	Healthcheck string `json:"health_check"`
	Services	ServiceMap `json:"services"`
}

type Service struct{
	Group        string          `json:"group"`
	Path         string          `json:"path"`
	Timeout      int             `json:"timeout"`
	Data         interface{}     `json:"data"`
}

type GilmourTopic string

type ServiceMap map[GilmourTopic]Service

func health_check(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("OK"))
}


func main()  {
	r := mux.NewRouter()
	r.HandleFunc("/health_check",health_check)
	log.Println("Test service running ...")
	req := Request{}
	req.Port = "8081"
	req.Healthcheck = "/health_check"
	m := make(map[GilmourTopic]Service)
	m["test"]=Service{"test","/test",10,""}
	req.Services = m

	fmt.Println(req)
	sendData, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return
	}
	//fmt.Println(string(sendData))

	byteStream := bytes.NewBuffer(sendData)
	//fmt.Println(byteStream)
	resp ,err := http.Post("http://127.0.0.1:8080/nodes","application/json",byteStream)

	if err != nil {
		log.Println("Error: ", err)
		return
	}

	defer resp.Body.Close()

	body ,_ := ioutil.ReadAll(resp.Body)
	http.ListenAndServe(":8081",r)
	fmt.Printf("Response body : %s",body)


}

