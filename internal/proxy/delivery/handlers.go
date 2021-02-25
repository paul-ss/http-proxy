package delivery

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Delivery struct {}

func (d *Delivery) Proxy(c *gin.Context) {

	req, err := http.NewRequest(c.Request.Method, c.Request.URL.String(), nil)
	if err != nil {
		log.Println("delivery, Proxy: " + err.Error())
		return
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if resp == nil || err != nil {
		log.Println("delivery, Proxy: resp == nil or err != nil")
		return
	}


	//if err := resp.Header.Write(c.Writer); err != nil {
	//	log.Println("delivery, Proxy: " + err.Error())
	//}

	c.Writer.WriteHeader(resp.StatusCode)
	c.

	//defer resp.Body.Close()
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//c.Writer.Write(body)
}