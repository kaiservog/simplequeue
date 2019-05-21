package main

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func createToken() (string, error) {
	secret := os.Getenv("SQ-SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"time": time.Now().Unix()})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func login(p string) bool {
	ok := os.Getenv("SHA1-PWD") //"GuOqAxqkaL7E2Hr1LUb8PjLX7dE=" change to list of pwds implement ':'
	hasher := sha1.New()
	hasher.Write([]byte(p))

	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha == ok
}

func validade(ts string) error {

	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secret := []byte(os.Getenv("SQ-SECRET"))
		return secret, nil
	})

	if err != nil {
		return err
	}

	if token.Valid {
		return nil
	}

	return errors.New("invalid token")
}
