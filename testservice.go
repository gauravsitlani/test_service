package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Request struct {
	Port        string     `json:"port"`
	Healthcheck string     `json:"health_check"`
	Services    ServiceMap `json:"services"`
}

type Service struct {
	Group   string      `json:"group"`
	Path    string      `json:"path"`
	Timeout int         `json:"timeout"`
	Data    interface{} `json:"data"`
}

type CreateNodeResponse struct {
	ID          string `json:"id"`
	PublishPort string `json:"publish_port"`
	Status      int    `json:"status"`
}

type GilmourTopic string

var myID string

type ServiceMap map[GilmourTopic]Service

func health_check(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func sendRequest() {
	req := Request{}
	req.Port = "8081"
	req.Healthcheck = "/health_check"
	m := make(map[GilmourTopic]Service)
	m["test"] = Service{"test", "/test", 10, ""}
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
	resp, err := http.Post("http://127.0.0.1:8080/nodes", "application/json", byteStream)

	if err != nil {
		log.Println("Error: ", err)
		return
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	responseData := CreateNodeResponse{}
	//====================================Tepmorary fix for "404 Not Found" appended at the end of the json file===========================
	body = body[:len(body)-19]
	//=====================================================================================================================================
	err = json.Unmarshal(body, &responseData)

	if err != nil {
		log.Print(err)
		log.Printf("%s", body)
		return
	}
	log.Print("Response body : %v", responseData)
	myID = responseData.ID
}

func getMyService() {
	requester := "http://127.0.0.1:8080/nodes/" + myID
	resp, err := http.Get(requester)
	log.Println(requester)
	if err != nil {
		log.Print(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Print(string(body))

}

func DeleteMe() {
	requester := "http://127.0.0.1:8080/nodes/" + myID
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", requester, nil)
	if err != nil {
		log.Print("Request creation error: ", err)
		return
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Print("Request error: ", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Print(string(body))
}

func FuncRunner() {
	sendRequest()
	<-time.After(time.Second * 8)
	getMyService()
	<-time.After(time.Second * 8)
	DeleteMe()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health_check", health_check)
	log.Println("Test service running ...")
	go FuncRunner()
	http.ListenAndServe(":8081", r)

}
