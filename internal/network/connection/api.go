package connection

import (
	"fmt"
	"github.com/paul-ss/http-proxy/internal/network/router"
	"net"
	"net/http"
	"os"
)

type Api struct {
	conn net.Conn
	router *router.Router
}


func NewApi(conn net.Conn, router *router.Router) *Api{
	return &Api{
		conn: conn,
		router: router,
	}
}

func (a *Api) Handle(r *http.Request) {
	defer a.conn.Close()

	fmt.Println("URI: " + r.RequestURI)
	resp := a.router.GetResponse(r)

	resp.Write(os.Stdout)


	if err := resp.Write(a.conn); err != nil {
		fmt.Println("Api-Handle: " + err.Error())
	}
}
