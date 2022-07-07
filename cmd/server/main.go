package main

import (
	"fmt"
	"time"

	"3mard.github.com/ably/pkg/server"
)

func main() {
	fmt.Println("Start server...")

	srv := server.NewServer(":8080", time.Duration(100)*time.Second)

	err := srv.Start()
	if err != nil {
		fmt.Println("Error:", err)
	}

}
