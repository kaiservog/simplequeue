package main

import (
	"testing"
)

func TestPut(t *testing.T) {
	q := newQ(2)

	for i := 0; i < 10; i++ {
		q.put("teste")
	}
}

/*
func TestPut(t *testing.T) {
	q := newQ(2)
	q.put("0")
	q.put("1")

	if "0" != q.get() {
		t.Fatal("first element must be 0")
	}
	if "1" != q.get() {
		t.Fatal("first element must be 1")
	}
}

func TestPutOverflow(t *testing.T) {
	q := newQ(2)
	q.put("0")
	q.put("1")
	q.put("2")

	if "1" != q.get() {
		t.Fatal("first element must be 2")
	}
	if "2" != q.get() {
		t.Fatal("first element must be 1")
	}
}

func TestPutLifetime(t *testing.T) {
	q := newQ(2)
	q.put("0")
	q.put("1")
	q.get()
	q.get()
	q.put("2")
	q.put("3")
	q.get()
	q.put("4")

	if "3" != q.get() {
		t.Fatal("first element must be 2")
	}
}
*/
