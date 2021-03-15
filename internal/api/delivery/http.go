package delivery

import (
	"fmt"
	"github.com/paul-ss/http-proxy/internal/api"
	"github.com/paul-ss/http-proxy/internal/api/usecase"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Delivery struct {
	uc api.IUsecase
}

func NewDelivery() *Delivery {
	d := &Delivery{}
	uc := usecase.NewUsecase(d)
	d.uc = uc

	return d
}

func (d *Delivery) Send(req *http.Request) (*http.Response, error) {
	cli := http.Client{}
	return cli.Do(req)
}


func (d *Delivery) GetRequests(w http.ResponseWriter, r *http.Request) {
	reqs, err := d.uc.GetRequests()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	for _, rr := range reqs {
		_, err := w.Write([]byte(fmt.Sprintf("%5d %s %s", rr.Id, rr.Method, rr.Path)))
		if err != nil {
			log.Println("Delivery-GetRequests-Write: " + err.Error())
			w.WriteHeader(500)
			return
		}
	}

	w.WriteHeader(200)
}

func (d *Delivery) GetRequestById(w http.ResponseWriter, r *http.Request) {
	id, err := getId("requests", r.URL.Path)
	if err != nil {
		log.Println("Delivery-GetRequestById-getId: " + err.Error())
		return
	}

	req, err := d.uc.GetRequestById(id)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, err = w.Write([]byte(req.Req))
	if err != nil {
		log.Println("Delivery-GetRequestById-Write: " + err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func (d *Delivery) RepeatById(w http.ResponseWriter, r *http.Request) {
	id, err := getId("repeat", r.URL.Path)
	if err != nil {
		log.Println("Delivery-RepeatById-getId: " + err.Error())
		return
	}

	req, err := d.uc.RepeatById(id)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, err = w.Write(req)
	if err != nil {
		log.Println("Delivery-RepeatById-Write: " + err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func (d *Delivery) ScanById(w http.ResponseWriter, r *http.Request) {
	id, err := getId("scan", r.URL.Path)
	if err != nil {
		log.Println("Delivery-ScanById-getId: " + err.Error())
		return
	}

	req, err := d.uc.ScanById(id)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, err = w.Write(req)
	if err != nil {
		log.Println("Delivery-ScanById-Write: " + err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func getId(prefix, path string) (int32, error) {
	idS := strings.TrimPrefix(strings.Trim(path, "/"), prefix + "/")
	id, err :=strconv.Atoi(idS)
	return int32(id), err
}
