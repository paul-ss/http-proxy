package network

import (
	"context"
	config "github.com/paul-ss/http-proxy/configs"
	"github.com/paul-ss/http-proxy/internal/api/delivery"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type ApiServer struct {
	srv *http.Server
}

func NewApiServer() *ApiServer {
	mux := http.NewServeMux()
	d := delivery.NewDelivery()

	mux.HandleFunc("/requests", d.GetRequests)
	mux.HandleFunc("/requests/", d.GetRequestById)
	mux.HandleFunc("/repeat/", d.RepeatById)
	mux.HandleFunc("/scan/", d.ScanById)


	return &ApiServer{
		srv: &http.Server{
			Addr: config.C.ApiAddress,
			Handler: mux,
		},
	}
}

func (s *ApiServer) Run() {
	go func() {
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ApiServer ListenAndServe: %v", err)
		}
	}()

	log.Println("ApiServer running at " + config.C.ApiAddress)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down ApiServer...")

	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("ApiServer Shutdown: %v", err)
		return
	}

	log.Printf("ApiServer stopped")
}


