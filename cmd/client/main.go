package main

import (
	"flag"
	"fmt"
	"log"

	"3mard.github.com/ably/pkg/client"
	"3mard.github.com/ably/pkg/hash"
)

func main() {
	address := flag.String("address", "", "address of the server localhost:8080")
	flag.Parse()
	if *address == "" {
		fmt.Println("address is required")
		return
	}
	cl := client.NewClient(*address)
	err := ConnectToServerWithRetrial(cl)
	if err != nil {
		log.Fatalf("error connecting to server %v", err)
	}
	err = cl.Handshake(client.HandshakeWithNumberOfMessages(10))
	if err != nil {
		panic(err)
	}
	checksumMessage, err := cl.ReadChecksum()
	if err != nil {
		panic(err)
	}
	fmt.Println(checksumMessage)
	data := []int32{}
	for i := 0; i < 3; i++ {
		number, err := cl.ReadNumber()
		if err != nil {
			panic(err)
		}
		data = append(data, number)
	}
	err = cl.Disconnect()
	if err != nil {
		panic(err)
	}

	err = cl.Connect()
	if err != nil {
		panic(err)
	}

	err = cl.Handshake(client.HandshakeWithOffset(3))
	if err != nil {
		panic(err)
	}

	for i := 0; i < 7; i++ {
		number, err := cl.ReadNumber()
		if err != nil {
			panic(err)
		}
		fmt.Println("Server sent number:", number)
		data = append(data, number)
	}
	fmt.Println("Client received data:", hash.CalculateChecksum(data))

}

func ConnectToServerWithRetrial(cl *client.Client) error {
	return nil
}
