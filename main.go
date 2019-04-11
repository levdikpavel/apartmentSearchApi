package main

import (
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

	log.Println("Appartment Api Service started at", GConfig.ServiceUrl)
	err = http.ListenAndServe(GConfig.ServiceUrl, router)
	if err != nil {
		log.Println(err)
	}
}

func search(w http.ResponseWriter, request *http.Request) {

}

func add(w http.ResponseWriter, request *http.Request) {
	data, _ := ioutil.ReadAll(request.Body)

	w.WriteHeader(200)
	w.Write(data)
}
