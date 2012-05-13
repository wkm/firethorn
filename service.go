package main

import (
	"log"
	"math/rand"
	"tumblr/redis"
)

type Service struct {
	retries int           // maximum retries
	pool    []redis.Redis // connection pool
}

// increment the given key
func (s *Service) Incr(key string) error {
	_, err := s.pickConn().Incr(key)
	return err
}

// decrement the given key
func (s *Service) Decr(key string) error {
	_, err := s.pickConn().Decr(key)
	return err
}

// reads the value of the key from one node and scales by the pool size
func (s *Service) Read(key string) int64 {
	val, err := s.pickConn().GetInt(key)

	// FIXME this should retry against another connection

	if err != nil {
		return val
	}
	log.Printf("could not read from Redis connection: %s", err)
	return 0
}

// randomly chooses a connection from the pool
func (s *Service) pickConn() *redis.Redis {
	redisId := rand.Intn(len(s.pool))
	return &s.pool[redisId]
}
