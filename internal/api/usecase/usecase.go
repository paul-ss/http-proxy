package usecase

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/paul-ss/http-proxy/internal/api"
	"github.com/paul-ss/http-proxy/internal/api/repository"
	"github.com/paul-ss/http-proxy/internal/domain"
	"log"
	"net/http"
	"os"
	"strings"
)

type Usecase struct {
	repo api.IRepository
	net api.INetwork
}

func NewUsecase(n api.INetwork) *Usecase {
	return &Usecase{
		repo: repository.NewDatabase(),
		net: n,
	}
}

func (uc *Usecase) StoreRequest(r http.Request) error {
	buff := bytes.NewBuffer([]byte{})
	if err := r.Write(buff); err != nil {
		log.Println("UC-StoreRequest-Write: " + err.Error())
		return err
	}

	rReq := domain.StoreRequest{
		Method: r.Method,
		Host: strings.Split(r.Host, ":")[0],
		Path: r.URL.String(),
		Req: buff.String(),
	}

	if _, err := uc.repo.StoreRequest(&rReq); err != nil {
		log.Println("UC-StoreRequest-repo: " + err.Error())
		return err
	}

	return nil
}


func (uc *Usecase) GetRequests() ([]domain.RequestShort, error) {
	req, err := uc.repo.GetShortRequests()
	if err != nil {
		log.Println("UC-GetRequests-repo: " + err.Error())
	}

	return req, err
}

func (uc *Usecase) GetRequestById(id int32) (*domain.Request, error) {
	req, err := uc.repo.GetRequestById(id)
	if err != nil {
		log.Println("UC-GetRequests-repo: " + err.Error())
	}

	return req, err
}

func (uc *Usecase) RepeatById(id int32) ([]byte, error) {
	req, err := uc.repo.GetRequestById(id)
	if err != nil {
		log.Println("UC-RepeatById-repo: " + err.Error())
		return nil, err
	}

	hReq, err := readRequest(req)
	if err != nil {
		log.Println("UC-RepeatById-ReadReq: " + err.Error())
		return nil, err
	}

	resp, err := uc.net.Send(hReq)
	if err != nil {
		log.Println("UC-RepeatById-Send: " + err.Error())
		return nil, err
	}

	respBuf := bytes.NewBuffer([]byte{})
	if err := resp.Write(respBuf); err != nil {
		log.Println("UC-RepeatById-Write: " + err.Error())
		return nil, err
	}

	return respBuf.Bytes(), nil
}

func (uc *Usecase) ScanById(id int32) ([]byte, error) {
	req, err := uc.repo.GetRequestById(id)
	if err != nil {
		log.Println("UC-ScanById-repo: " + err.Error())
		return nil, err
	}

	hReq, err := readRequest(req)
	if err != nil {
		log.Println("UC-ScanById-ReadReq: " + err.Error())
		return nil, err
	}

	res := bytes.NewBuffer([]byte{})
	strs, err := mapFile()
	if err != nil {
		log.Println("UC-ScanById-map: " + err.Error())
		return nil, err
	}
	log.Println("Scan started")

	for i, s := range strs {
		if i % 50 == 0 {
			log.Printf("%d routes done\n", i)
		}
		cloneReq := hReq.Clone(context.Background())
		cloneReq.URL.Path = s

		resp, err := uc.net.Send(cloneReq)
		if err != nil {
			log.Println("UC-ScanById-Send: " + err.Error())
			continue
		}

		if resp.StatusCode == 404 {
			continue
		}

		res.WriteString(fmt.Sprintf("%d - /%s\n", resp.StatusCode, s))
	}

	return res.Bytes(), nil
}





type ProxyUsecase struct {
	repo api.IRepository
}

func NewProxyUsecase() *ProxyUsecase {
	return &ProxyUsecase{
		repo: repository.NewDatabase(),
	}
}

func (uc *ProxyUsecase) StoreRequest(r http.Request) error {
	buff := bytes.NewBuffer([]byte{})
	if err := r.Write(buff); err != nil {
		log.Println("UC-StoreRequest-Write: " + err.Error())
		return err
	}

	rReq := domain.StoreRequest{
		Method: r.Method,
		Host: strings.Split(r.Host, ":")[0],
		Path: r.URL.String(),
		Req: buff.String(),
	}

	if _, err := uc.repo.StoreRequest(&rReq); err != nil {
		log.Println("UC-StoreRequest-repo: " + err.Error())
		return err
	}

	return nil
}


func readRequest(r *domain.Request) (*http.Request, error) {
	hReq, err := http.ReadRequest(bufio.NewReader(bytes.NewBufferString(r.Req)))
	if err != nil {
		log.Println("UC-ScanById-ReadReq: " + err.Error())
		return nil, err
	}
	hReq.RequestURI = ""
	hReq.URL.Scheme = "http"
	hReq.URL.Host = r.Host

	return hReq, nil
}

func mapFile() ([]string, error) {
	file, err := os.Open("configs/dicc.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var res []string
	s := bufio.NewScanner(file)
	for s.Scan() {
		res = append(res, s.Text())
	}

	return res, nil
}