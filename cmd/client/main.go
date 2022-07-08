package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"3mard.github.com/ably/pkg/client"
	"3mard.github.com/ably/pkg/hash"
)

const (
	MaxConnectRetry = uint64(3)
)

func main() {
	address := flag.String("address", "", "address of the server localhost:8080")
	flag.Parse()
	if *address == "" {
		fmt.Println("address is required")
		return
	}
	cl := client.NewClient(*address)
	err := cl.ConnectWithRetrial(MaxConnectRetry)
	if err != nil {
		log.Fatalf("error connecting to server %v", err)
	}
	err = cl.Handshake(client.HandshakeWithNumberOfMessages(10))
	if err != nil {
		log.Fatalf("error handshaking with server %v", err)
	}
	checksumMessage, err := cl.ReadChecksum()
	if err != nil {
		log.Fatalf("error reading checksum %v", err)
	}
	log.Printf("serve checksum: %v", checksumMessage.Checksum)
	data := []int32{}
	progress := 0
	for progress < checksumMessage.NumberOfMessages-1 {
		message, err := cl.ReadSequence()
		if progress >= checksumMessage.NumberOfMessages {
			break
		}
		// connection dropped
		if err == io.EOF {
			err = cl.ConnectWithRetrial(MaxConnectRetry)
			if err != nil {
				log.Fatalf("error connecting to server %v", err)
			}
			err = cl.Handshake(client.HandshakeWithOffset(progress))
			if err != nil {
				log.Fatalf("error handshaking with server %v", err)
			}
			continue
		}
		if err != nil {
			log.Fatalf("error reading number %v", err)
		}
		data = append(data, message.Sequence)
		progress = message.Index
	}
	log.Printf("Calculated hash %v", hash.CalculateChecksum(data))
}
