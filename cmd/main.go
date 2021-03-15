package main

import (
	"github.com/paul-ss/http-proxy/internal/database"
	"github.com/paul-ss/http-proxy/internal/network"
	"log"
)

func main()  {
	defer database.Close()

	proxy := network.NewProxyServer()
	go proxy.Run()

	api := network.NewApiServer()
	api.Run()



	//quit := make(chan os.Signal)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//
	//<-quit
	//log.Println("Shutting down server...")

	proxy.Stop()

	log.Println("ProxyServer exiting")
}
