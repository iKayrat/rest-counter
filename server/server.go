package server

import (
	"github.com/gin-gonic/gin"
	"github.com/iKayrat/rest-counter/db"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(redis *db.Store) *Server {

	server := &Server{store: *redis}
	router := gin.Default()

	router.POST("rest/counter/add/:i", server.addCounter)
	router.POST("rest/counter/sub/:i", server.subCounter)
	router.GET("rest/counter/val", server.getCounter)

	router.POST("rest/substr/find", server.getSubstr)
	router.POST("rest/self/find/:str", server.findStr)

	router.POST("rest/hash/calc", server.genHash)
	router.GET("rest/hash/result/:id", server.getHash)

	router.POST("rest/email/check/", server.chekEmail)
	router.POST("rest/inn/check/", server.checkINN)

	server.router = router

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
