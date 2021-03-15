package network

import (
	"bufio"
	config "github.com/paul-ss/http-proxy/configs"
	"github.com/paul-ss/http-proxy/internal/network/cert"
	"github.com/paul-ss/http-proxy/internal/network/connection"
	"github.com/paul-ss/http-proxy/internal/network/router"
	"log"
	"net"
	"net/http"
	"sync"
)

type ProxyServer struct {
	listener net.Listener
	quit     chan interface{}
	wg 		 sync.WaitGroup
	certs	 *cert.Certs
	router   *router.Router
}

func NewProxyServer() *ProxyServer {
	ln, err := net.Listen("tcp", config.C.ProxyAddress)
	if err != nil {
		log.Fatal("network-NewProxyServer-listen: " + err.Error())
	}

	return &ProxyServer{
		listener: ln,
		quit:     make(chan interface{}),
		certs:    cert.NewCerts(),
		router:   router.NewRouter(),
	}
}

func (s *ProxyServer) Run() {
	log.Println("ProxyServer running at " + config.C.ProxyAddress)

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

func (s *ProxyServer) Stop() {
	log.Println("Shutting down ProxyServer...")

	close(s.quit)
	s.listener.Close()
	s.wg.Wait()

	log.Printf("ProxyServer stopped")
}




func (s *ProxyServer) handleConnection(conn net.Conn) {
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		conn.Close()
		log.Println("network-handleConnection: " + err.Error())
		return
	}

	req.Header.Get("Proxy-Connection")


	if len(req.Header.Get("Proxy-Connection")) > 0 {
		if req.Method == "CONNECT" {
			c := connection.NewHttpsConn(conn, s.certs)
			c.Handle(req)
		} else {
			c := connection.NewHttpConn(conn)
			c.Handle(req)
		}
	} else {
		connection.NewApi(conn, s.router).Handle(req)
	}
}



