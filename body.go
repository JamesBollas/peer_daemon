package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"context"
	"bytes"
	"errors"
)

var background = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

type ClientConnection struct {
	Socket string
	SocketType string
	Service string
}

type ClientConnectionManager struct {
	connections map[ClientConnection]struct{}
}

func (manager ClientConnectionManager) createConnection(service string) ClientConnection{
	connection := ClientConnection{"./test.sock", "unix", service}
	manager.connections[connection] = struct{}{}
	return connection
}

func (manager ClientConnectionManager) serviceConnections(service string) []ClientConnection {
	connections := make([]ClientConnection, 0)
	for connection, _ := range manager.connections{
		if connection.Service == service{
			connections = append(connections, connection)
		}
	}
	return connections
}

var manager = ClientConnectionManager{map[ClientConnection]struct{}{}}

func HandleMessage(body []byte, remoteAddress string) error {
	service, message, isValid := getService(body)
	if !isValid{
		return errors.New("message is not valid")
	}
	connections := manager.serviceConnections(service)
	for _, connection := range connections {
		sendMessageToClient(connection, message)
	}

    val, err := rdb.Get(background, "key").Result()
    if err != nil {
        panic(err)
    }
	fmt.Println("key", val)
	return nil
}

func getService(body []byte) (string, []byte, bool) {
	split := bytes.SplitN(body, []byte("\n"), 2)
	if len(split) != 2{
		return "", nil, false
	}
	return string(split[0]), split[1], true
}

func CreateConnection(service string) ClientConnection {
	return manager.createConnection(service)
}

func sendMessageToClient(connection ClientConnection, message []byte){
	fmt.Println(message)
}