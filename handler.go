package main

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

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
func DeleteQ(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	err := validade(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	id := token + "." + params["name"]

	err = helper.deleteQ(id)

	if err != nil {
		if err.Error() == "no queue" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		io.WriteString(w, err.Error())
		return
	}
}

func CreateQ(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	err := validade(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	id := token + "." + params["name"]

	err = helper.createQ(10, id)
	if err != nil && err.Error() == "queue exists, DELETE it" {
		w.WriteHeader(http.StatusConflict)
		io.WriteString(w, err.Error())
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func doQueueLock(token, name string) *sync.Mutex {
	qmux := processing[token+"."+name]

	if qmux == nil {
		qmux = &sync.Mutex{}
		qmux.Lock()
		processing[token+"."+name] = qmux
	} else {
		qmux.Lock()
	}

	return qmux

}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	err := validade(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	name := params["name"]
	qmux := doQueueLock(token, name)
	defer qmux.Unlock()

	id := token + "." + name

	alg := newCircularAlg(helper)

	e, err := alg.get(id)

	if err != nil {
		if err.Error() == "empty" {
			w.WriteHeader(http.StatusNoContent)
		} else if err.Error() == "no queue" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		}
		return
	}

	io.WriteString(w, e)
}

func PutMessage(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	err := validade(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	name := params["name"]
	id := token + "." + name
	qmux := doQueueLock(token, name)
	defer qmux.Unlock()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	c := buf.String()

	alg := newCircularAlg(helper)
	err = alg.put(c, id)

	if err != nil {
		if err.Error() == "no queue" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
