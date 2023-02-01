package api

import (
	"e2e-test/src/cache"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	paramKey   = "key"
	paramValue = "value"

	routeKey      = "/:" + paramKey
	routeKeyValue = "/:" + paramKey + "/:" + paramValue
)

func (s *server) HandleGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName(paramKey)
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)

		return
	}

	value, err := s.Redis.Get(key)
	if errors.Is(err, cache.ErrNotFound) {
		http.Error(w, "key not found", http.StatusNotFound)

		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get value for key: %s", err), http.StatusInternalServerError)

		return
	}

	_, err = w.Write([]byte(value))
	if err != nil {
		log.Printf("failed to write response: %s", err)
	}
}

func (s *server) HandleSet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName(paramKey)
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)

		return
	}

	value := ps.ByName(paramValue)
	if value == "" {
		http.Error(w, "missing value", http.StatusBadRequest)

		return
	}

	if err := s.Redis.Set(key, value); err != nil {
		http.Error(w, fmt.Sprintf("failed to set value for key: %s", err), http.StatusInternalServerError)

		return
	}

	log.Printf("key '%s' populated with value '%s'", key, value)

	w.WriteHeader(http.StatusCreated)
}

func (s *server) HandleDel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName(paramKey)
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)

		return
	}

	if err := s.Redis.Del(key); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete key: %s", err), http.StatusInternalServerError)

		return
	}

	log.Printf("key '%s' deleted", key)

	w.WriteHeader(http.StatusOK)
}
