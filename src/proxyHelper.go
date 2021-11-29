package main

import(
	"golang.org/x/net/proxy"
	"net/http"
	//"net/url"
	"bytes"
	"io/ioutil"
	"os"
	"fmt"
	"strings"
	"os/exec"
)

func PostThroughProxy(address string, message []byte, headers map[string]string) (int, []byte, error){
	torLocal := os.Getenv("PROXY_SOCKET")
	torSocketType := os.Getenv("PROXY_SOCKET_TYPE")
	fmt.Println(torSocketType)
	tbDialer, err := proxy.SOCKS5(torSocketType,torLocal ,nil, proxy.Direct)
	if err != nil{
		return 500, nil, err
	}
	
	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	client := &http.Client{Transport: tbTransport}

	messageReader := bytes.NewReader(message)

	req, err := http.NewRequest("POST", address , messageReader)
	if headers != nil{
		for key, value := range headers {
			key = strings.TrimRight(key, "\r\n")
			value = strings.TrimRight(value, "\r\n")
			req.Header.Add(key, value)
		}
	}
	if err != nil{
		panic(err)
	}
	fmt.Println("got past headers")
	//req.Header.Add("User-Agent", "myClient")

	resp, err := client.Do(req)
	if err != nil{
		fmt.Println(err)
		return 500, nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, body, nil
}

func StartProxy(){
	proxyExecutable := os.Getenv("PROXY_EXECUTABLE")
	proxyConfig := os.Getenv("PROXY_CONFIG")
	cmd := exec.Command(proxyExecutable, "-f", proxyConfig)
	go cmd.Run()
}