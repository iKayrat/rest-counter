package main

import (
	"log"

	"github.com/iKayrat/rest-counter/db"
	"github.com/iKayrat/rest-counter/server"
)

var (
	listenAddress = "127.0.0.1:8080"
	redisAddress  = "localhost:6379"
)

func main() {
	//new redis conn
	store, err := db.NewStore(redisAddress)
	if err != nil {
		log.Fatal("cannot connect to redis: ", err)
	}

	//new server
	server := server.NewServer(store)
	if err != nil {
		log.Fatal("new server err:", err)
	}

	err = server.Start(listenAddress)
	if err != nil {
		log.Fatal("server run() err:", err)
	}
}
