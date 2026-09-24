package main

import "github.com/yeabsraayehualem/nehmya_saas/cmd/server"


func main(){
	s := server.NewServer()
	s.Run()
}