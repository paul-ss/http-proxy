package connection

import (
	"github.com/paul-ss/http-proxy/internal/network/http"
)

type Connection interface {
	Handle(r *http.Request)
}
