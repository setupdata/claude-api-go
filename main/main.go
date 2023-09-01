package main

import (
	"fmt"
	"github.com/imroc/req/v3"
)

func main() {

	client := req.DevMode()
	client.ImpersonateChrome()
	client.SetProxyURL("http://127.0.0.1:7890")

	testUrl := "https://tls.peet.ws/api/clean"

	resp, err := client.R().Get(testUrl)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = resp.Body.Close() }()
}
