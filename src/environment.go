package main

import (
	"github.com/joho/godotenv"
)

func LoadEnvironment(){
	err := godotenv.Load(EnvironmentPath)
	if err != nil{
		//todo set default values here?
		panic("no environment")
	}
}