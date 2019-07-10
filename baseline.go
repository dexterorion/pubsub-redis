package main

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

// Service service.
type Service struct {
	pool *redis.Pool
	conn redis.Conn
}

// New return new service
func New(input *NewInput) *Service {
	if input == nil {
		log.Fatal("input is required")
	}
	var redispool *redis.Pool
	redispool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", input.RedisURL)
		},
	}

	// Get a connection
	conn := redispool.Get()
	defer conn.Close()
	// Test the connection
	_, err := conn.Do("PING")
	if err != nil {
		log.Fatalf("can't connect to the redis database, got error:\n%v", err)
	}

	return &Service{
		pool: redispool,
		conn: conn,
	}
}

// NewInput input for constructor.
type NewInput struct {
	RedisURL string
}

// Subscribe subscribe
func (s *Service) Subscribe(key string, msg chan []byte) error {
	rc := s.pool.Get()
	psc := redis.PubSubConn{Conn: rc}
	if err := psc.PSubscribe(key); err != nil {
		return err
	}

	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.PMessage:
				msg <- v.Data
			}
		}
	}()
	return nil
}

// Publish publish key value
func (s *Service) Publish(key string, value string) error {
	conn := s.pool.Get()
	conn.Do("PUBLISH", key, value)
	return nil
}
