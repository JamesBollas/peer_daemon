package main

import (
	"net/http"
	//"log"
	"fmt"
	"io/ioutil"
	"io"
	"yukon_go/encryptionHelper"
	"yukon_go/torHelper"
)

func main() {
	exposed := http.NewServeMux()

	exposed.HandleFunc("/", indexRoute)
	exposed.HandleFunc("/hashcookie", hashCookieRoute)

	locals := http.NewServeMux()
	locals.HandleFunc("/sendmessage", sendMessage)

	fmt.Printf("Starting server at port 8080\n")
	go http.ListenAndServe(":8081", locals)
	http.ListenAndServe(":8080", exposed)
}

func indexRoute(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	remoteAddress := request.Header.Get("remoteAddress")
	signature := request.Header.Get("signature")
	hashType := request.Header.Get("hash")
	verification := encryptionHelper.VerifySignature(remoteAddress, hashType, signature, body)
	if !verification{
		writer.WriteHeader(403)
		return
	}


	fmt.Println(body)

	//fmt.Println(user)
	//fmt.Fprintf(writer, body)
}

func hashCookieRoute(writer http.ResponseWriter, request * http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	signature := encryptionHelper.SignHash(body)
	io.WriteString(writer, signature)
}

func sendMessage(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	address := request.Header.Get("address")
	signature := encryptionHelper.SignHash(body)
	myAddress := "http://wmjrfxz2ikhfi2vmm7jgtttktmjivibbcxsgzupc5m55px76a2ihdzad.onion"
	
	headers := map[string]string {"remoteAddress":myAddress,"signature":signature,"hash":"sha256"}

	torHelper.PostWithHeader(address, body, headers)

}