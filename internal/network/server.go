package network

import (
	"fmt"
	config "github.com/paul-ss/http-proxy/configs"
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
			c := NewConnection(conn)
			c.handleConnection()
			s.wg.Done()
		}()
	}
}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}



type Connection struct {
	ClientConn net.Conn
	ServerConn net.Conn
	srvClosed bool
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		ClientConn: conn,
		srvClosed: true,
	}
}


func (c *Connection) handleConnection() {
	for {
		req := NewRequest()
		if err := req.Parse(c.ClientConn); err != nil {
			if err == io.EOF {
				log.Println("EOF found")
				return
			}

			log.Println("network-handleConnection: " + err.Error())
		}

		fmt.Println(string(req.Bytes()))

		var resp *Response
		var err error

		if len(req.Url.Host) > 0 {
			if req.Method == "CONNECT" {
				resp, err = c.handleHttpsProxy(req)
			} else {
				c.handleHttpProxy(req)
				return
			}
		} else {
			//resp, err = c.handleLocal(req)
		}

		if err != nil {
			log.Println("handleConnection error: " + err.Error())
			return
		}

		if _, err := c.ClientConn.Write(resp.Bytes()); err != nil {
			if err == io.EOF {
				log.Println("EOF found")
				return
			}

			log.Println("network-handleConnection: " + err.Error())
		}
	}

}


func (c *Connection) handleHttpProxy(req *Request) (*Response, error) {
	host, ok := req.Headers["Host"]
	if !ok {
		log.Println("Can't find host header")
		return nil, fmt.Errorf("can't find host header")
	}

	conn, err := net.Dial("tcp", host + ":80")
	if err != nil {
		log.Println("Can't connect to host: " + err.Error())
		return nil, err
	}
	defer conn.Close()

	req.Url.Scheme = ""
	req.Url.Host = ""
	delete(req.Headers, "Proxy-Connection")

	if _, err := conn.Write(req.Bytes()); err != nil {
		log.Println("Error writing to server: " + err.Error())
		return nil, err
	}

	resp := NewResponse()
	if err := resp.Parse(conn); err != nil {
		log.Println("Error parse resp: " + err.Error())
		return nil, err
	}

	return resp, nil
}


func (c *Connection) handleHttpsProxy(req *Request) {

}

