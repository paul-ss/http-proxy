package connection

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

type HttpConn struct {
	ClientConn net.Conn
	ServerConn net.Conn
}


func NewHttpConn(conn net.Conn) *HttpConn {
	return &HttpConn{
		ClientConn: conn,
	}
}

func (c *HttpConn) OpenServerConn(r *http.Request) error {
	conn, err := net.Dial("tcp", r.Host + ":80")
	if err != nil {
		log.Println("Can't connect to host: " + err.Error())
		return err
	}

	c.ServerConn = conn
	return nil
}

func (c *HttpConn) Handle(r *http.Request) {
	defer c.ClientConn.Close()

	if err := c.OpenServerConn(r); err != nil {
		log.Println("HttpConn-Handle: " + err.Error())
		return
	}
	defer c.ServerConn.Close()

	//r.Url.Scheme = ""
	//r.Url.Host = ""
	//delete(r.Headers, "Proxy-Connection")
	r.Header.Del("Proxy-Connection")

	if err := r.Write(c.ServerConn); err != nil {
		log.Println("HttpConn-Handle: Error writing to server: " + err.Error())
		return
	}

	resp, err := http.ReadResponse(bufio.NewReader(c.ServerConn), r)
	if  err != nil {
		log.Println("HttpConn-Handle: Error parse resp: " + err.Error())
		return
	}

	if err := resp.Write(c.ClientConn); err != nil {
		log.Println("HttpConn-Handle: Error writing to client: " + err.Error())
		return
	}
}


