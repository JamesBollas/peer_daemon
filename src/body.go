package main

import (
	"bytes"
	"errors"
)

// check for format validity, then save to database
func HandleMessage(body []byte, remoteAddress string) error {
	service, message, isValid := getService(body)
	if !isValid{
		return errors.New("message is not valid")
	}

	userIdentifier, idType := getUserIdentifier(remoteAddress)

	AddMessage(service, userIdentifier, idType, message)

	return nil
}

// gets utf-8 formatted service string from message. All messages must have format: {utf-8 service string}{newline}{binary message}
// service type may enforce additional restrictions on the message binary format, these are not handled here.
func getService(body []byte) (string, []byte, bool) {
	split := bytes.SplitN(body, []byte("\n"), 2)
	if len(split) != 2{
		return "", nil, false
	}
	return string(split[0]), split[1], true
}

// TODO, will return username associated with remoteAddress if one exists, first check if cached, then hit user server.
// If no username exists return remote address as unique identifier
func getUserIdentifier(remoteAddress string) (string, string){
	return remoteAddress, "address"
}
