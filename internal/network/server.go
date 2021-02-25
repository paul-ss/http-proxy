package network

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	config "github.com/paul-ss/http-proxy/configs"
	"github.com/paul-ss/http-proxy/internal/proxy/delivery"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	HttpSrv *http.Server
}

func NewServer() Server {
	r := gin.Default()

	proxyHandler := delivery.Delivery{}

	r.Any("/", proxyHandler.Proxy)

	return Server{
		HttpSrv: &http.Server{
			Addr: config.C.ProxyAddress,
			Handler: r,
		},
	}
}

func (s *Server) Run() {
	go func() {
		if err := s.HttpSrv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	log.Printf("listen on: %s\n", config.C.ProxyAddress)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")


	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.HttpSrv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
