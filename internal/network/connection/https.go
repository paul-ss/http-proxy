package connection

import (
	"fmt"
	"github.com/paul-ss/http-proxy/internal/network/http"
	"log"
	"net"
)

type HttpsConn struct {
	ClientConn net.Conn
	ServerConn net.Conn
}


func NewHttpsConn(conn net.Conn) *HttpsConn {
	return &HttpsConn{
		ClientConn: conn,
	}
}

func (c *HttpsConn) OpenServerConn(r *http.Request) error {
	host, ok := r.Headers["Host"]
	if !ok {
		log.Println("Can't find host header")
		return fmt.Errorf("can't find host header")
	}

	conn, err := net.Dial("tcp", host+":443")
	if err != nil {
		log.Println("Can't connect to host: " + err.Error())
		return err
	}

	c.ServerConn = conn
	return nil
}

func (c *HttpsConn) Handle(r *http.Request) {
	defer c.ClientConn.Close()

	if err := c.OpenServerConn(r); err != nil {
		log.Println("HttpsConn-Handle: " + err.Error())
		return
	}
	defer c.ServerConn.Close()

	r.Url.Scheme = ""
	r.Url.Host = ""
	delete(r.Headers, "Proxy-Connection")

	if _, err := c.ServerConn.Write(r.Bytes()); err != nil {
		log.Println("HttpsConn-Handle: Error writing to server: " + err.Error())
		return
	}

	resp := http.NewResponse()
	if err := resp.Parse(c.ServerConn); err != nil {
		log.Println("HttpsConn-Handle: Error parse resp: " + err.Error())
		return
	}

	if _, err := c.ClientConn.Write(resp.Bytes()); err != nil {
		log.Println("HttpsConn-Handle: Error writing to client: " + err.Error())
		return
	}
}
