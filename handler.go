package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func GetToken(w http.ResponseWriter, r *http.Request) {
	log.Println("request get token")
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
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, token)
}

func DeleteQ(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	err := validade(token)
	h := tokenToHash(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("Authorization fail")
		return
	}

	params := mux.Vars(r)
	id := h + "." + params["name"]
	log.Println("request delete q", id)

	err = helper.deleteQ(id)

	if err != nil {
		if err.Error() == "no queue" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Println(err)
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

	h := tokenToHash(token)
	params := mux.Vars(r)
	id := h + "." + params["name"]
	log.Println("request create q", id)

	err = helper.createQ(10, id)
	if err != nil && err.Error() == "queue exists, DELETE it" {
		w.WriteHeader(http.StatusConflict)
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	err := validade(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	h := tokenToHash(token)
	params := mux.Vars(r)
	name := params["name"]

	id := h + "." + name
	log.Println("request get message", id)
	qmux := doQueueLock(id, name)
	defer qmux.Unlock()

	alg := newCircularAlg(helper)

	e, err := alg.get(id)

	if err != nil {
		if err.Error() == "empty" {
			w.WriteHeader(http.StatusNoContent)
		} else if err.Error() == "no queue" {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, err.Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		}
		log.Println(err)
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

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	c := buf.String()

	if len(c) > 256 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "max message size is 256")

	}

	h := tokenToHash(token)
	params := mux.Vars(r)
	name := params["name"]
	id := h + "." + name

	log.Println("request put message", id)
	qmux := doQueueLock(id, name)
	defer qmux.Unlock()

	alg := newCircularAlg(helper)
	err = alg.put(c, id)

	if err != nil {
		if err.Error() == "no queue" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		}
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func tokenToHash(token string) string {
	hasher := sha1.New()
	hasher.Write([]byte(token))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return hash[:len(hash)-1]
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
