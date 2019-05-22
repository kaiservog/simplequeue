package main

type queue struct {
	id    string
	idxI  int
	idxO  int
	depth int
}

type repository interface {
	queue(id string) (*queue, error)
	getMessage(idx int, id string) (string, error)
	deleteMessage(idx int, id string) error
	putMessage(idx int, message, id string) error
	updateIdx(idxI, idxO int, id string) error
	createQ(depth int, id string) error
}
