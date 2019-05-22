package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

var processing map[string]*sync.Mutex
var helper *redisHelper

func main() {
	processing = make(map[string]*sync.Mutex)

	address := "192.168.1.109" //TODO
	password := ""             //TODO
	port := "6379"
	h, err := newRedisHelper(address+":"+port, password)
	helper = h

	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/token", GetToken).Methods("GET")
	router.HandleFunc("/q/{name}", CreateQ).Methods("POST")
	router.HandleFunc("/q/{name}", DeleteQ).Methods("DELETE")
	router.HandleFunc("/q/{name}", GetMessage).Methods("GET")
	router.HandleFunc("/q/{name}", PutMessage).Methods("PUT")

	log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
}
