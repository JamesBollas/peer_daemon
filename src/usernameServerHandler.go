package main

import(
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"os"
)

func GetAddressFromUsername(username string) string{
	fmt.Println("getting address from username")
	fmt.Println(username)

	client := &http.Client{}

	messageReader := bytes.NewReader([]byte(username))

	req, err := http.NewRequest("POST", os.Getenv("UNS_ADDRESS") , messageReader)
	resp, err := client.Do(req)
	if err != nil{
		fmt.Println("something went wrong with the uns request!")
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}