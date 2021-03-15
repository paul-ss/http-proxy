package main

import (
	"github.com/paul-ss/http-proxy/internal/database"
	"github.com/paul-ss/http-proxy/internal/network"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	defer database.Close()

	srv := network.NewServer()
	go srv.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	srv.Stop()

	log.Println("Server exiting")
}
