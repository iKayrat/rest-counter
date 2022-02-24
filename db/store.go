package db

import (
	"log"

	"github.com/go-redis/redis"
)

type Store struct {
	Client *redis.Client
}

func NewStore(address string) (*Store, error) {

	cl := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	ping, err := cl.Ping().Result()
	if err != nil {
		return nil, err
	}
	log.Println("ping: ", ping)

	return &Store{Client: cl}, nil
}
