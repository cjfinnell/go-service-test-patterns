package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type server struct {
	redis  RedisWrapper
	router http.Handler
}

func newServer(conf *config) *server {
	s := server{
		redis: newRedisWrap(conf),
	}
	router := httprouter.New()

	router.GET("/status", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
	})

	router.GET(routeKey, authWrap(s.handleGet))
	router.POST(routeKeyValue, authWrap(s.handleSet))
	router.DELETE(routeKey, authWrap(s.handleDel))

	s.router = requestLogger(router)

	return &s
}

func (s *server) run() error {
	addr := ":8080"

	defer s.redis.Close()

	log.Printf("server listening on %s", addr)

	return http.ListenAndServe(addr, s.router)
}
