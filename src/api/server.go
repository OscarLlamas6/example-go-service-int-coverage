package api

import (
	"e2e-test/src/cache"
	"e2e-test/src/config"
	"e2e-test/src/utils"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type server struct {
	Redis  cache.RedisWrapper
	Router http.Handler
}

func NewServer(conf *config.Config) *server {
	s := server{
		Redis: cache.NewRedisWrap(conf),
	}
	router := httprouter.New()
	router.GET(routeKey, s.HandleGet)
	router.POST(routeKeyValue, s.HandleSet)
	router.DELETE(routeKey, s.HandleDel)
	s.Router = utils.RequestLogger(utils.AuthWrap(router))

	return &s
}

func (s *server) Run() error {
	addr := ":8080"

	defer s.Redis.Close()

	log.Printf("server listening on %s", addr)

	return http.ListenAndServe(addr, s.Router)
}
