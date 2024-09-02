package main

import (
	"GoExamGatewayAPI/internal/config"
	"GoExamGatewayAPI/internal/server"
	"GoExamGatewayAPI/internal/stopsignal"
	"log"
)

func main() {
	cfg := config.MustLoad()

	srv := server.New(cfg)
	srv.Start()
	log.Printf("Server started\n")

	stopsignal.Stop()

	srv.Shutdown()
}
