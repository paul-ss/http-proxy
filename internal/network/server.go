package network

import (
	"fmt"
	config "github.com/paul-ss/http-proxy/configs"
	"github.com/paul-ss/http-proxy/internal/network/cert"
	"github.com/paul-ss/http-proxy/internal/network/connection"
	"github.com/paul-ss/http-proxy/internal/network/http"
	"io"
	"log"
	"net"
	"sync"
)

type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg 		 sync.WaitGroup
	certs	 *cert.Certs
}

func NewServer() *Server {
	ln, err := net.Listen("tcp", config.C.ProxyAddress)
	if err != nil {
		log.Fatal("network-NewServer-listen: " + err.Error())
	}

	return &Server{
		listener: ln,
		quit:     make(chan interface{}),
		certs:    cert.NewCerts(),
	}
}

func (s *Server) Run() {
	log.Println("Server running at " + config.C.ProxyAddress)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
				case <-s.quit:
					return
				default:
					log.Println("network-Run-Accept: " + err.Error())
			}
		}

		s.wg.Add(1)
		go func() {
			s.handleConnection(conn)
			s.wg.Done()
		}()
	}
}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}




func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close() //!!

	req := http.NewRequest()
	if err := req.Parse(conn); err != nil {
		if err == io.EOF {
			log.Println("EOF found")
			return
		}

		log.Println("network-handleConnection: " + err.Error())
		return
	}

	fmt.Println(string(req.Bytes()))


	if _, ok := req.Headers["Proxy-Connection"]; ok {
		if req.Method == "CONNECT" {
			c := connection.NewHttpsConn(conn, s.certs)
			c.Handle(req)
		} else {
			c := connection.NewHttpConn(conn)
			c.Handle(req)
		}
	} else {
		//resp, err = c.handleLocal(req)
		log.Println("No handlers run")
		log.Printf("Host: %s, Path: %s, Scheme: %s", req.Url.Host, req.Url.Path, req.Url.Scheme)
	}
}



