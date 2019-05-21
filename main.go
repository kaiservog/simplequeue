package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var qq map[string]*Q

func main() {
	qq = make(map[string]*Q)

	router := mux.NewRouter()
	router.HandleFunc("/token", GetToken).Methods("GET")
	router.HandleFunc("/q/{id}", CreateQ).Methods("POST")
	router.HandleFunc("/q/{id}", GetElm).Methods("GET")
	router.HandleFunc("/q/{id}", PutQ).Methods("PUT")
	//router.HandleFunc("/q/{id}", DeletePerson).Methods("DELETE")

	log.Fatal(http.ListenAndServe(os.Getenv("SQ-PORT"), router))
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	pwd := r.Header.Get("Authorization")
	ok := login(pwd)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "login failure")
		return
	}

	token, err := createToken() //must register token in a db for blacklist
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, token)
}

func CreateQ(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	err := validade(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	q := newQ(10)

	qq[params["id"]] = q
	w.WriteHeader(http.StatusOK)
}

func GetElm(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	err := validade(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	q := qq[params["id"]]
	if q == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	q.mux.Lock()
	defer q.mux.Unlock()

	e, err := q.get()

	if err != nil {
		if err.Error() == "empty" {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

	}

	io.WriteString(w, e)
	return
}

func PutQ(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	err := validade(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	q := qq[params["id"]]

	if q == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	q.mux.Lock()
	defer q.mux.Unlock()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	c := buf.String()

	q.put(c)
	w.WriteHeader(http.StatusOK)
}
