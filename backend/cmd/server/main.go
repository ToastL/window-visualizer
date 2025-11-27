package main

import (
	"backend/internal/config"
	"backend/internal/websocket"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
		os.Exit(0)
	}()

	http.HandleFunc("/ws", websocket.WsHandler)

	go websocket.BroadcastLoop(ctx)

	addr := "127.0.0.1:" + config.AppConfig.Port
	log.Printf("WebSocket server on ws://%s/ws", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
