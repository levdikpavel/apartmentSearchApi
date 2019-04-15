package main

import (
	"encoding/json"
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
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading body"))
		return
	}
	var req ApartmentSearchRequest
	err = parseRequest(data, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := Manager.searchApartments(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
}

func add(w http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading body"))
		return
	}
	var req Apartment
	err = parseRequest(data, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if req.CorpusName == "" {
		req.CorpusName = "default"
	}

	result, err := Manager.addApartment(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
}

func parseRequest(data []byte, req interface{}) error {
	err := json.Unmarshal(data, req)
	if err != nil {
		log.Printf("Error while parsing request. %v", err)
	}
	return err
}
