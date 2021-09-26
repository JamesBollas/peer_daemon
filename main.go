package main

import (
	"net/http"
	//"log"
	"fmt"
	"io/ioutil"
	"os"
	"net"
	//"io"
	//"yukon_go/authenticationHelper"
	//"yukon_go/torHelper"
)

func (p program) run() {
	//run by the service wrapper. essentially the main function.
	InitializeKeys()
	LoadEnvironment()


	exposed := http.NewServeMux()
	exposed.HandleFunc("/", indexRoute)
	exposed.HandleFunc("/publickey", publicKeyRoute)

	locals := http.NewServeMux()
	locals.HandleFunc("/sendmessage", sendMessage)

	localListener, _ := net.Listen("tcp",os.Getenv("LOCAL_SOCKET"))
	hiddenServiceListener, _ := net.Listen("unix",os.Getenv("HIDDEN_SERVICE_SOCKET"))

	fmt.Printf("Starting server\n")
	go http.Serve(localListener, locals)
	http.Serve(hiddenServiceListener, exposed)
}

func indexRoute(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("recieved a thingy")
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

	
	headers := map[string]string {"remoteAddress":myAddress,"signature":signature,"cryptoStandard":"ed25519"}

	PostThroughProxy(address, body, headers)

}

func getMyAddress() string{
	hostname, err := ioutil.ReadFile(os.Getenv("HOSTNAME_PATH"))
	if err != nil{
		panic("no valid hostname file at location")
	}
	return string(hostname)
}
