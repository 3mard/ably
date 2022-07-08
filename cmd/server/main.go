package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"3mard.github.com/ably/pkg/server"
)

func main() {
	address := flag.String("address", "", "address of the server localhost:8080")
	flag.Parse()
	if *address == "" {
		fmt.Println("address is required")
		return
	}
	log.Println("Server started on ", *address)
	srv := server.NewServer(*address, time.Duration(100)*time.Second)
	err := srv.Start()
	if err != nil {
		log.Fatal(err)
	}

}
