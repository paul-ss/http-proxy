package connection

import (
	"crypto/tls"
	"fmt"
	"github.com/paul-ss/http-proxy/internal/network/cert"
	"github.com/paul-ss/http-proxy/internal/network/http"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type HttpsConn struct {
	ClientConn net.Conn
	ServerConn net.Conn
	certs      cert.ICerts
	wg 		   sync.WaitGroup
}


func NewHttpsConn(conn net.Conn, certs cert.ICerts) *HttpsConn {
	return &HttpsConn{
		ClientConn: conn,
		certs: certs,
	}
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
	srvConn, err := tls.Dial(
		"tcp",
		strings.TrimLeft(clientReq.Url.String(), "/"),
		&tls.Config{})
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


	cer, err := c.certs.GetCert(clientReq.Url.Host)
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

	//if err := c.OpenServerConn(r); err != nil {
	//	log.Println("HttpsConn-Handle: " + err.Error())
	//	return
	//}
	//defer c.ServerConn.Close()

	if err := c.Connect(r); err != nil {
		log.Println("HttpsConn-Handle-Connect: " + err.Error())
		return
	}

	//c.ClientConn.SetDeadline(time.Now().Add(5 *time.Second))
	//c.ServerConn.SetDeadline(time.Now().Add(5 *time.Second))

	//f := func(dst io.WriteCloser, src io.ReadCloser) {
	//	//defer dst.Close()
	//	//defer src.Close()
	//
	//	var err error
	//	for {
	//		if _, err = io.Copy(dst, src); err != nil {
	//			log.Println("f: " + err.Error())
	//			break
	//		}
	//	}
	//	wg.Done()
	//}

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

		n, err := io.Copy(c.ServerConn, c.ClientConn)
		if err != nil {
			log.Println("HandleClientToSrv: " + err.Error())
			break
		}
		if n == 0 {
			log.Println("HandleClientToSrv: 0" )
			break
		}


		//b := make([]byte, 200)
		//n, err := c.ClientConn.Read(b)
		//if err != nil {
		//	log.Println("HandleClientToSrv-Read: " + err.Error())
		//	break
		//}
		//
		//_, err = fmt.Fprint(c.ServerConn, b[:n])
		//if err != nil {
		//	log.Println("HandleClientToSrv-Write: " + err.Error())
		//	break
		//}

	}
	c.wg.Done()
}

func (c *HttpsConn) HandleSrvToClient() {
	for {
		if err := c.ServerConn.SetDeadline(time.Now().Add(10 *time.Second)); err != nil {
			log.Println("HandleSrvToClient: " + err.Error())
			break
		}

		n, err := io.Copy(c.ClientConn, c.ServerConn)
		if err != nil {
			log.Println("HandleSrvToClient: " + err.Error())
			break
		}
		if n == 0 {
			log.Println("HandleClientToSrv: 0" )
			break
		}
	}
	c.wg.Done()
}