package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var GConfig *Config
var Manager *MysqlManager

func main() {
	GConfig = LoadConfig()
	Manager = new(MysqlManager)
	err := Manager.Connect()
	if err != nil {
		log.Fatal("DbManager connect failed.\n", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/search", search)
	router.HandleFunc("/add", add)

	log.Println("Apartment Api Service started at", GConfig.ServiceUrl)
	err = http.ListenAndServe(GConfig.ServiceUrl, router)
	if err != nil {
		log.Println(err)
	}
}

func search(w http.ResponseWriter, request *http.Request) {

}

func add(w http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading body"))
		return
	}
	req, err := parseRequest(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Print(req)

	w.WriteHeader(200)
	w.Write(data)
}

func parseRequest(data []byte) (Apartment, error) {
	var req Apartment
	err := json.Unmarshal(data, &req)
	if err != nil {
		log.Printf("Error while parsing request. %v", err)
	}
	return req, err
}
