package cert

import (
	"crypto/tls"
	"fmt"
	config "github.com/paul-ss/http-proxy/configs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type ICerts interface {
	GetCert(host string) (*tls.Certificate, error)
}

type Certs struct {
	certs map[string]*tls.Certificate
	rmList []string
}

func NewCerts() *Certs {
	return &Certs{
		certs: make(map[string]*tls.Certificate),
	}
}

func (c *Certs) addCert(host string, cert *tls.Certificate) {
	c.certs[host] = cert
	c.rmList = append(c.rmList, host)

	if len(c.rmList) > config.C.MaxInMemoryCerts {
		delete(c.certs, c.rmList[0])
		c.rmList = c.rmList[1:]
		fmt.Println("Deleted cert for host " + host)
	}

	fmt.Println("Added cert for host " + host)
}

func (c *Certs) openCert(host string) (*tls.Certificate, error) {
	if _, err := os.Stat("cert.key"); os.IsNotExist(err) { // TODO
		log.Println("Cert: Cert key doesn't exist")
		return nil, fmt.Errorf("cert: Cert key doesn't exist")
	}

	out, err := exec.Command("./scripts/gen-certs.sh", strings.Split(host, ":")[0]).Output()
	if err != nil {
		log.Println("Cert: Can't gen cert: " + err.Error())
		return nil, err
	}

	fmt.Println(string(out))

	key, err := ioutil.ReadFile("cert.key")
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(out, key)
	if err != nil {
		return nil, err
	}

	c.addCert(host, &cert)
	return &cert, nil
}

func (c *Certs) GetCert(host string) (*tls.Certificate, error) {
	host = strings.Split(host, ":")[0]

	cert, ok := c.certs[host]
	if ok {
		return cert, nil
	}

	cert, err := c.openCert(host)
	if err != nil {
		log.Println("GetCert: " + err.Error())
		return nil, err
	}

	return cert, nil
}
