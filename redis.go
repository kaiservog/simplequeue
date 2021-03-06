package main

import (
	"errors"
	"strconv"

	"github.com/go-redis/redis"
)

type redisHelper struct {
	client *redis.Client
}

const (
	_idxi    = ".idxi"
	_idxo    = ".idxo"
	_message = ".m."
	_depth   = ".d"
)

func newRedisHelper(url string) (*redisHelper, error) {
	opt, err := redis.ParseURL(url)

	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	_, err = client.Ping().Result()

	if err != nil {
		return nil, err
	}

	r := redisHelper{
		client: client}

	return &r, nil
}

func (r *redisHelper) queue(id string) (*queue, error) {
	i, err := r.idxI(id)

	if err == redis.Nil {
		return nil, errors.New("no queue")
	}

	if err != nil {
		return nil, err
	}

	o, err := r.idxO(id)
	if err == redis.Nil {
		return nil, errors.New("no queue")
	}

	if err != nil {
		return nil, err
	}

	d, err := r.depth(id)
	if err == redis.Nil {
		return nil, errors.New("no queue")
	}

	if err != nil {
		return nil, err
	}

	q := queue{
		id:    id,
		idxI:  i,
		idxO:  o,
		depth: d}

	return &q, nil
}

func (r *redisHelper) createQ(depth int, id string) error {
	q, err := r.queue(id)

	if q != nil {
		return errors.New("queue exists, DELETE it")
	}

	if err != nil && err.Error() != "no queue" {
		return err
	}

	d := strconv.Itoa(depth)

	err = r.client.Set(id+_depth, d, 0).Err()
	if err != nil {
		return err
	}

	err = r.client.Set(id+_idxi, 0, 0).Err()
	if err != nil {
		return err
	}

	err = r.client.Set(id+_idxo, 0, 0).Err()
	return err
}

func (r *redisHelper) deleteQ(id string) error {
	q, err := r.queue(id)

	if q == nil {
		return errors.New("no queue")
	}

	_, err = r.client.Del(id + _depth).Result()
	if err != nil {
		return err
	}

	_, err = r.client.Del(id + _idxi).Result()
	if err != nil {
		return err
	}

	_, err = r.client.Del(id + _idxo).Result()
	return err
}

func (r *redisHelper) idxI(id string) (int, error) {
	return r.client.Get(id + _idxi).Int()
}

func (r *redisHelper) idxO(id string) (int, error) {
	return r.client.Get(id + _idxo).Int()
}

func (r *redisHelper) getMessage(idx int, id string) (string, error) {
	idxStr := strconv.Itoa(idx)
	m, err := r.client.Get(id + _message + idxStr).Result()

	if err == redis.Nil {
		return "", errors.New("empty")
	}

	if err != nil {
		return "", err
	}

	return m, nil
}

func (r *redisHelper) deleteMessage(idx int, id string) error {
	idxStr := strconv.Itoa(idx)
	_, err := r.client.Del(id + _message + idxStr).Result()

	return err
}

func (r *redisHelper) putMessage(idx int, message, id string) error {
	idxStr := strconv.Itoa(idx)
	key := id + _message + idxStr
	err := r.client.Set(key, message, 0).Err()

	return err
}

func (r *redisHelper) updateIdx(idxI, idxO int, id string) error {

	if idxI != -1 {
		idxiStr := strconv.Itoa(idxI)
		err := r.client.Set(id+_idxi, idxiStr, 0).Err()

		if err != nil {
			return err
		}
	}

	if idxO != -1 {
		idxoStr := strconv.Itoa(idxO)
		err := r.client.Set(id+_idxo, idxoStr, 0).Err()

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *redisHelper) depth(id string) (int, error) {
	return r.client.Get(id + _depth).Int()

}
