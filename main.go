package main

import (
	"log"
	"net"
)

func main() {
	cfg, err := loadConfig("config.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ln, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s", cfg.Listen)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go handleConn(conn, cfg)
	}
}
