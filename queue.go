package main

import (
	"errors"
	"fmt"
	"sync"
)

/* Q is the queue struct */
type Q struct {
	idxRead int64
	idxPut  int64

	depth int64
	ee    []string

	mux sync.Mutex
}

const nilkey = "xxx"

func newQ(depth int64) *Q {
	q := Q{
		idxRead: 0,
		idxPut:  0,
		depth:   depth}

	q.ee = make([]string, depth)
	var i int64
	for i = 0; i < depth; i++ {
		q.ee[i] = nilkey
	}
	return &q
}

func incrementOrReset(idx, max int64) int64 {
	idx++

	if idx >= max {
		return 0
	}

	return idx
}

func (q *Q) put(elm string) {
	fmt.Println("oops")
	q.ee[q.idxPut] = elm

	if q.idxPut == q.idxRead {
		if q.depth < q.idxRead && q.ee[q.idxRead+1] != nilkey {
			q.idxRead = incrementOrReset(q.idxRead, q.depth)
		}
	}

	q.idxPut = incrementOrReset(q.idxPut, q.depth)
}

func (q *Q) get() (string, error) {
	e := q.ee[q.idxRead]

	if e == nilkey {
		return "", errors.New("empty")
	}

	q.ee[q.idxRead] = nilkey
	q.idxRead++

	if q.idxRead >= q.depth {
		q.idxRead = 0
	}

	return e, nil
}
