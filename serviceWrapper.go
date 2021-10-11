package main

import (
   "fmt"
   "github.com/kardianos/service"
   "os"
)

const serviceName = "peer_daemon"
const serviceDescription = "allows peer connections through internet proxy with global referents (tor)"

type program struct{}

func (p program) Start(s service.Service) error {
   fmt.Println(s.String() + " started")
   go p.run()
   return nil
}

func (p program) Stop(s service.Service) error {
   fmt.Println(s.String() + " stopped")
   os.Remove(os.Getenv("HIDDEN_SERVICE_SOCKET"))
   os.Remove(os.Getenv("LOCAL_SOCKET"))
   os.Remove(os.Getenv("PROXY_SOCKET"))
   return nil
}

func main() {
   serviceConfig := &service.Config{
      Name:        serviceName,
      DisplayName: serviceName,
      Description: serviceDescription,
   }
   prg := &program{}
   s, err := service.New(prg, serviceConfig)
   if err != nil {
      fmt.Println("Cannot create the service: " + err.Error())
   }
   err = s.Run()
   if err != nil {
      fmt.Println("Cannot start the service: " + err.Error())
   }
}