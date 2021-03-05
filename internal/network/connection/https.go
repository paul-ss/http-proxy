package connection

import (
	"crypto/tls"
	"fmt"
	"github.com/paul-ss/http-proxy/internal/network/http"
	"io"
	"log"
	"net"
	"os/exec"
	"sync"
	"time"
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

func (c *HttpsConn) GenCert(r *http.Request) error {
	if err := exec.Command("./scripts/gen-serts.sh", r.Url.String()).Run(); err != nil {
		log.Println("Can't gen cert: " + err.Error())
		return err
	}

	return nil
}

func (c *HttpsConn) OpenServerConn(r *http.Request) error {
	host, ok := r.Headers["Host"]
	if !ok {
		log.Println("Can't find host header")
		return fmt.Errorf("can't find host header")
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Println("Can't connect to host: " + err.Error())
		return err
	}

	c.ServerConn = conn
	return nil
}

func (c *HttpsConn) Connect(clientReq *http.Request) error {
	srvConn, err := tls.Dial("tcp", clientReq.Url.String(), &tls.Config{})
	if err != nil {
		log.Println("HttpsConn-Connect-Dial: " + err.Error())
		return err
	}
	c.ServerConn = srvConn

	clResp := http.NewResponse()
	clResp.Status = 200
	clResp.Protocol = clientReq.Protocol
	clResp.Message = "Connection Established"
	clResp.Headers["Proxy-agent"] = "paul-s-proxy"

	if _, err := c.ClientConn.Write(clResp.Bytes()); err != nil {
		log.Println("HttpsConn-Connect: Error writing to server: " + err.Error())
		return err
	}

	if err := c.GenCert(clientReq); err != nil {
		log.Println("HttpsConn-Connect-GenCert: " + err.Error())
		return err
	}

	cer, err := tls.LoadX509KeyPair(
		"certs/www.google.com.crt",
		"cert.key",
		)
	if err != nil {
		log.Println("HttpsConn-Connect-Load: " + err.Error())
		return err
	}

	clConn := tls.Server(c.ClientConn, &tls.Config{
		Certificates: []tls.Certificate{cer},
	})

	c.ClientConn = clConn
	return nil
}

func (c *HttpsConn) Handle(r *http.Request) {
	defer c.ClientConn.Close()
	wg := sync.WaitGroup{}

	if err := c.OpenServerConn(r); err != nil {
		log.Println("HttpsConn-Handle: " + err.Error())
		return
	}
	defer c.ServerConn.Close()

	if err := c.Connect(r); err != nil {
		log.Println("HttpsConn-Handle-Connect: " + err.Error())
		return
	}

	c.ClientConn.SetDeadline(time.Now().Add(5 *time.Second))
	c.ServerConn.SetDeadline(time.Now().Add(5 *time.Second))

	f := func(dst io.WriteCloser, src io.ReadCloser) {
		//defer dst.Close()
		//defer src.Close()

		var err error
		for {
			if _, err = io.Copy(dst, src); err != nil {
				log.Println("f: " + err.Error())
				break
			}
		}
		wg.Done()
	}

	wg.Add(2)
	go f(c.ClientConn, c.ServerConn)
	go f(c.ServerConn, c.ClientConn)
	wg.Wait()
}
