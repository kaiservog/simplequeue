package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

const (
	redisURL = "REDIS_URL"
)

var processing map[string]*sync.Mutex
var helper *redisHelper

func main() {
	log.Println("Starting...")
	processing = make(map[string]*sync.Mutex)
	h, err := newRedisHelper(os.Getenv(redisURL))
	helper = h
	if err != nil {
		panic(err)
	}

	log.Println("Redis up")

	router := mux.NewRouter()
	router.HandleFunc("/token", GetToken).Methods("GET")
	router.HandleFunc("/q/{name}", CreateQ).Methods("POST")
	router.HandleFunc("/q/{name}", DeleteQ).Methods("DELETE")
	router.HandleFunc("/q/{name}", GetMessage).Methods("GET")
	router.HandleFunc("/q/{name}", PutMessage).Methods("PUT")

	log.Println("Mux up")
	log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
}
