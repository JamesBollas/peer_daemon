package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
	"net"
)

func (p program) run() {
	//run by the service wrapper. essentially the main function.
	InitializeKeys()
	LoadEnvironment()
	StartProxy()

	//routes visible as hidden service
	exposed := http.NewServeMux()
	exposed.HandleFunc("/", indexRoute)
	exposed.HandleFunc("/publickey", publicKeyRoute)

	//routes visible on local machine/network only
	locals := http.NewServeMux()
	locals.HandleFunc("/sendmessage", sendMessage)
	locals.HandleFunc("/getmessage", getMessage)
	locals.HandleFunc("/getmessageids", getMessageIds)

	localListener, _ := net.Listen("unix",os.Getenv("LOCAL_SOCKET"))
	hiddenServiceListener, _ := net.Listen("unix",os.Getenv("HIDDEN_SERVICE_SOCKET"))

	fmt.Printf("Starting server\n")
	go http.Serve(localListener, locals)
	http.Serve(hiddenServiceListener, exposed)
}

// receiving messages from remote
func indexRoute(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)

	remoteAddress := request.Header.Get("remoteAddress")
	signature := request.Header.Get("signature")
	cryptoStandard := request.Header.Get("cryptoStandard")
	myAddress := getMyAddress()

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

// sending public key to remote
func publicKeyRoute(writer http.ResponseWriter, request * http.Request) {
	key := []byte(MyPublicKey())
	writer.Write(key)
}

// hit by local to send message to index route of other user/device
func sendMessage(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	address := request.Header.Get("address")

	signature := SignBody(body, address)
	myAddress := getMyAddress()
	headers := map[string]string {"remoteAddress":myAddress,"signature":signature,"cryptoStandard":"ed25519"}

	statusCode, postReturn, _ := PostThroughProxy(address, body, headers)
	writer.WriteHeader(statusCode)
	writer.Write(postReturn)
}

// hit by local to retrieve a message by id from the sqlite database
func getMessage(writer http.ResponseWriter, request *http.Request) {
	id := request.Header.Get("id")
	message, _ := GetMessage(id)
	writer.Write(message)
}

// hit by local to retrieve all message ids from sqlite database
func getMessageIds(writer http.ResponseWriter, request *http.Request) {
	messageIds := GetMessageIds()
	for _, messageId := range messageIds{
		writer.Write([]byte(messageId + "\n"))
	}
}

// not a route, returns the hidden service address of this machine
func getMyAddress() string{
	hostname, err := ioutil.ReadFile(os.Getenv("HOSTNAME_PATH"))
	if err != nil{
		fmt.Println("cannot access hostname file at location")
		fmt.Println(os.Getenv("HOSTNAME_PATH"))
		panic(err)
	}
	return "http://" + string(hostname)
}
