package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	//client := &http.Client{}


	req, err := http.NewRequest("GET", "http://google.com/123", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	b := bytes.NewBuffer([]byte{})
	fmt.Println(req.URL.String())
	req.Write(b)
	fmt.Println(b.String())


	//resp, err := client.Do(req)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	//b := bytes.NewBuffer([]byte{})
	//resp.Write(b)
	//fmt.Println(b.String())
}