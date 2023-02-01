package main

import (
	"e2e-test/src/api"
	"e2e-test/src/config"
	"log"
)

func main() {
	conf := config.NewConfig()

	srv := api.NewServer(conf)

	if err := srv.Run(); err != nil {
		log.Fatalf("server exited with error: %s", err)
	}
}
