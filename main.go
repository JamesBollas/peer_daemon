package main

import (
	"net/http"
	//"log"
	"fmt"
	"io/ioutil"
	"os"
	"net"
	"encoding/json"
	//"io"
	//"yukon_go/authenticationHelper"
	//"yukon_go/torHelper"
)

func (p program) run() {
	//run by the service wrapper. essentially the main function.
	InitializeKeys()
	LoadEnvironment()
	StartProxy()

	exposed := http.NewServeMux()
	exposed.HandleFunc("/", indexRoute)
	exposed.HandleFunc("/publickey", publicKeyRoute)

	locals := http.NewServeMux()
	locals.HandleFunc("/sendmessage", sendMessage)
	locals.HandleFunc("/connect", connectClient)

	localListener, _ := net.Listen("unix",os.Getenv("LOCAL_SOCKET"))
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
	// fmt.Println(remoteAddress)
	verification := VerifySignature(remoteAddress, myAddress, cryptoStandard, signature, body)
	if !verification{
		writer.WriteHeader(403)
		return
	}

	err := HandleMessage(body, remoteAddress)
	if err != nil{
		writer.WriteHeader(401)
		return
	}
	writer.WriteHeader(200)
	return
}

func publicKeyRoute(writer http.ResponseWriter, request * http.Request) {
	fmt.Println("sending public key")
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

	//fmt.Println(len(myAddress))
	//fmt.Println(len(signature))
	headers := map[string]string {"remoteAddress":myAddress,"signature":signature,"cryptoStandard":"ed25519"}
	//headers = map[string]string {"hi":"hi"}
	statusCode, postReturn, _ := PostThroughProxy(address, body, headers)
	writer.WriteHeader(statusCode)
	writer.Write(postReturn)
}

func connectClient(writer http.ResponseWriter, request *http.Request) {
	service := request.Header.Get("service")
	connection := CreateConnection(service)
	connectionJSON, _ := json.Marshal(&connection)
	writer.Write(connectionJSON)
}

func getMyAddress() string{
	hostname, err := ioutil.ReadFile(os.Getenv("HOSTNAME_PATH"))
	if err != nil{
		fmt.Println("cannot access hostname file at location")
		fmt.Println(os.Getenv("HOSTNAME_PATH"))
		panic(err)
	}
	return "http://" + string(hostname)
}
