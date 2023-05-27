package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

func MakeRequest() {
	url := "http://localhost/biz/get_users?user_id=2&message_id=2&auth_key=2&auth_id=ABCDE01234ABCDE0123793FdrQPTRkl897Etrw7T"
	method := "GET"

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http.DefaultClient.Timeout = 5 * time.Second
	client := &http.Client{
		Timeout: 10 * time.Second,
		// Transport: &http.Transport{
		// 	DialContext: (&net.Dialer{
		// 		Timeout:   5 * time.Second,
		// 		KeepAlive: 3 * time.Second,
		// 	}).DialContext,
		// 	ResponseHeaderTimeout: 5 * time.Second,
		// 	ExpectContinueTimeout: 1 * time.Second,
		// 	MaxIdleConnsPerHost:   -1,
		// },
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	wg.Done()
}

func main() {
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go MakeRequest()
	}
	wg.Wait()
}
