package main

import (
	"net/http"
	//"log"
	"fmt"
	"io/ioutil"
	//"io"
	//"yukon_go/authenticationHelper"
	//"yukon_go/torHelper"
)

func (p program) run() {
	//run by the service wrapper. essentially the main function.
	InitializeKeys()
	exposed := http.NewServeMux()

	exposed.HandleFunc("/", indexRoute)
	exposed.HandleFunc("/publickey", publicKeyRoute)

	locals := http.NewServeMux()
	locals.HandleFunc("/sendmessage", sendMessage)

	fmt.Printf("Starting server at port 8080\n")
	go http.ListenAndServe(":8081", locals)
	http.ListenAndServe(":8080", exposed)
}

func indexRoute(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	//fmt.Println(string(body))
	remoteAddress := request.Header.Get("remoteAddress")
	signature := request.Header.Get("signature")
	cryptoStandard := request.Header.Get("cryptoStandard")
	myAddress := getMyAddress()
	verification := VerifySignature(remoteAddress, myAddress, cryptoStandard, signature, body)
	if !verification{
		writer.WriteHeader(403)
		return
	}


	fmt.Println(body)
}

func publicKeyRoute(writer http.ResponseWriter, request * http.Request) {
	key := []byte(MyPublicKey())
	writer.Write(key)
}

func sendMessage(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	//fmt.Println(string(body))
	address := request.Header.Get("address")

	signature := SignBody(body, address)
	fmt.Println(address)
	fmt.Println(signature)
	myAddress := getMyAddress()

	
	headers := map[string]string {"remoteAddress":myAddress,"signature":signature,"hash":"sha256"}

	PostThroughProxy(address, body, headers)

}

func getMyAddress() string{
	return "http://wmjrfxz2ikhfi2vmm7jgtttktmjivibbcxsgzupc5m55px76a2ihdzad.onion"
}
