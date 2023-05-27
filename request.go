package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

func MakeRequest() {
	url := "http://localhost/biz/get_users?user_id=2&message_id=2&auth_key=2&auth_id=ABCDE01234ABCDE0123793FdrQPTRkl897Etrw7T"
	method := "GET"
  
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client {
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
	wg.Add(2000)
	for i := 0; i < 2000; i++ {
		go MakeRequest()
	}
	wg.Wait()
}
