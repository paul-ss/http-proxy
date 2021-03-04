package network

import (
	"fmt"
	config "github.com/paul-ss/http-proxy/configs"
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
}

func NewServer() *Server {
	ln, err := net.Listen("tcp", config.C.ProxyAddress)
	if err != nil {
		log.Fatal("network-NewServer-listen: " + err.Error())
	}

	return &Server{
		listener: ln,
		quit: make(chan interface{}),
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
	req := http.NewRequest()
	if err := req.Parse(conn); err != nil {
		if err == io.EOF {
			log.Println("EOF found")
			return
		}

		log.Println("network-handleConnection: " + err.Error())
	}

	fmt.Println(string(req.Bytes()))


	if len(req.Url.Host) > 0 {
		if req.Method == "CONNECT" {
			c := connection.NewHttpsConn(conn)
			c.Handle(req)
		} else {
			c := connection.NewHttpConn(conn)
			c.Handle(req)
			return
		}
	} else {
		//resp, err = c.handleLocal(req)
	}
}



