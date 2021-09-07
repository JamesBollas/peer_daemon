package main

import(
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"bytes"
	"io/ioutil"
)

func PostResponse (address string, message []byte ) ([]byte, error) {
	torLocal := "socks5://127.0.0.1:9050"

	tbProxyURL, _ := url.Parse(torLocal)

	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil{
		return nil, err
	}
	
	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	client := &http.Client{Transport: tbTransport}

	messageReader := bytes.NewReader(message)

	resp, err := client.Post(address, "bytes", messageReader)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func PostWithHeader(address string, message []byte, headers map[string]string) ([]byte, error){
	torLocal := "socks5://127.0.0.1:9050"

	tbProxyURL, _ := url.Parse(torLocal)

	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil{
		return nil, err
	}
	
	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	client := &http.Client{Transport: tbTransport}

	messageReader := bytes.NewReader(message)

	req, err := http.NewRequest("POST", address , messageReader)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	//req.Header.Add("User-Agent", "myClient")

	resp, err := client.Do(req)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}