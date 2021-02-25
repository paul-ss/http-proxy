package main

import "github.com/paul-ss/http-proxy/internal/network"

func main()  {
	srv := network.NewServer()
	srv.Run()
}
