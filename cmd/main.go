package main

import "nehmya/cmd/server"


func main(){
	s := server.NewServer()
	s.Run()
}