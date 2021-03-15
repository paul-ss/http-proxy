package api

import (
	"github.com/paul-ss/http-proxy/internal/domain"
	"net/http"
)

type IRepository interface {
	StoreRequest(req *domain.StoreRequest) (*domain.Request, error)
	GetShortRequests() ([]domain.RequestShort, error)
	GetRequestById(id int32) (*domain.Request, error)
}

type IUsecase interface {
	StoreRequest(r http.Request) error
	GetRequests() ([]domain.RequestShort, error)
	GetRequestById(id int32) (*domain.Request, error)
	RepeatById(id int32) ([]byte, error)
	ScanById(id int32) ([]byte, error)
}

type INetwork interface {
	Send(req *http.Request) (*http.Response, error)
}