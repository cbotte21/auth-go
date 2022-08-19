//TODO: Abstract api

package main

import (
	"log"
	"github.com/cbotte21/games-auth/service"
)

func main() {
	api, res := service.NewApi(5000)
	if !res || api.Start() != nil { //Start API Listener
		log.Fatal("Failed to initialize API.")
	}
}
