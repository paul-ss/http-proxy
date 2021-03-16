package connection

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/paul-ss/http-proxy/internal/api/usecase"
	"github.com/paul-ss/http-proxy/internal/network/cert"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type HttpsConn struct {
	ClientConn net.Conn
	ServerConn net.Conn
	certs      cert.ICerts
	wg 		   sync.WaitGroup
	uc 		   *usecase.ProxyUsecase
	host       string
}


func NewHttpsConn(conn net.Conn, certs cert.ICerts) *HttpsConn {
	return &HttpsConn{
		ClientConn: conn,
		certs: certs,
		uc: usecase.NewProxyUsecase(),
	}
}


func (c *HttpsConn) OpenServerConn(r *http.Request) error {
	conn, err := net.Dial("tcp", r.Host)
	if err != nil {
		log.Println("Can't connect to host: " + err.Error())
		return err
	}

	c.ServerConn = conn
	return nil
}

func (c *HttpsConn) Connect(r *http.Request) error {
	srvConn, err := tls.Dial("tcp", r.Host, &tls.Config{})
	if err != nil {
		log.Println("HttpsConn-Connect-Dial: " + err.Error())
		return err
	}
	c.ServerConn = srvConn

	clResp, err := http.ReadResponse(bufio.NewReader(bytes.NewBufferString(
		fmt.Sprintf(
			"HTTP/1.1 200 Connection Established\r\nProxy-agent: paul-s-proxy\r\n\r\n",
			))), nil)
	if err != nil {
		log.Println("HttpsConn-Connect-ReadResp: " + err.Error())
		return err
	}

	if err := clResp.Write(c.ClientConn); err != nil {
		log.Println("HttpsConn-Connect: Error writing to server: " + err.Error())
		return err
	}


	cer, err := c.certs.GetCert(r.Host)
	if err != nil {
		return err
	}

	clConn := tls.Server(c.ClientConn, &tls.Config{
		Certificates: []tls.Certificate{*cer},
	})

	c.ClientConn = clConn
	return nil
}

func (c *HttpsConn) Handle(r *http.Request) {
	defer c.ClientConn.Close()
	c.host = r.Host
	if err := c.Connect(r); err != nil {
		log.Println("HttpsConn-Handle-Connect: " + err.Error())
		return
	}
	defer c.ServerConn.Close()


	c.wg.Add(2)
	go c.HandleClientToSrv()
	go c.HandleSrvToClient()
	c.wg.Wait()
}

func (c *HttpsConn) HandleClientToSrv() {
	for {
		if err := c.ClientConn.SetDeadline(time.Now().Add(10 *time.Second)); err != nil {
			log.Println("HandleClientToSrv: " + err.Error())
			break
		}

		r := bufio.NewReader(c.ClientConn)
		req, err := http.ReadRequest(r)
		if err != nil {
			log.Println("HandleClientToSrv-Parse: " + err.Error())
			break
		}
		req.Host = c.host

		if err := c.uc.StoreRequest(*req); err != nil {
			log.Println("HandleClientToSrv-StoreRequest: " + err.Error())
			break
		}

		if err := req.Write(c.ServerConn); err != nil {
			log.Println("HandleClientToSrv-Write: " + err.Error())
			break
		}
	}
	c.wg.Done()
}

func (c *HttpsConn) HandleSrvToClient() {
	for {
		if err := c.ServerConn.SetDeadline(time.Now().Add(10 *time.Second)); err != nil {
			log.Println("HandleSrvToClient: " + err.Error())
			break
		}

		r := bufio.NewReader(c.ServerConn)
		resp, err := http.ReadResponse(r, nil)
		if err != nil {
			log.Println("HandleSrvToClient-Parse: " + err.Error())
			break
		}

		// store resp

		if err := resp.Write(c.ClientConn); err != nil {
			log.Println("HandleSrvToClient-Write: " + err.Error())
			break
		}
	}
	c.wg.Done()
}