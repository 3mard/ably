package tests

import (
	"testing"
	"time"

	cl "3mard.github.com/ably/pkg/client"
	"3mard.github.com/ably/pkg/hash"
	"3mard.github.com/ably/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	client := cl.NewClient(":8080")
	server := server.NewServer(":8080", 10*time.Second)
	go server.Start()
	err := client.Connect()
	assert.NoError(t, err)
}

func TestHandshake(t *testing.T) {
	client := cl.NewClient(":8080")
	server := server.NewServer(":8080", 10*time.Second)
	go server.Start()
	err := client.Connect()
	assert.NoError(t, err)
	err = client.Handshake(cl.HandshakeWithNumberOfMessages(10))
	assert.NoError(t, err)
}

func TestReceivingChecksum(t *testing.T) {
	client := cl.NewClient(":8080")
	server := server.NewServer(":8080", 10*time.Second)
	go server.Start()
	err := client.Connect()
	assert.NoError(t, err)
	err = client.Handshake(cl.HandshakeWithNumberOfMessages(10))
	assert.NoError(t, err)
	checksumMessage, err := client.ReadChecksum()
	assert.NoError(t, err)
	assert.NotNil(t, checksumMessage)
}

func TestErrorResumeReceivingNewClient(t *testing.T) {
	client := cl.NewClient(":8080", cl.WithClientId("test"))
	server := server.NewServer(":8080", 10*time.Second)
	go server.Start()
	err := client.Connect()
	assert.NoError(t, err)
	err = client.Handshake(cl.HandshakeWithNumberOfMessages(10), cl.HandshakeWithOffset(20))
	assert.NoError(t, err)
	errorMessage, err := client.ReadError()
	assert.NoError(t, err)
	assert.NotNil(t, errorMessage)
	assert.Equal(t, "Can't continue progress, please start over", errorMessage)
}

func TestReceivingSequence(t *testing.T) {
	client := cl.NewClient(":8080")
	server := server.NewServer(":8080", 10*time.Second)
	go server.Start()
	err := client.Connect()
	assert.NoError(t, err)
	err = client.Handshake(cl.HandshakeWithNumberOfMessages(10))
	assert.NoError(t, err)
	checksum, err := client.ReadChecksum()
	assert.NoError(t, err)
	data := make([]int32, 0)
	for i := 0; i < 10; i++ {
		number, err := client.ReadNumber()
		assert.NoError(t, err)
		data = append(data, number)
	}
	assert.Equal(t, checksum.Checksum, hash.CalculateChecksum(data))
}

func TestReceivingSequenceWithoutNumberOfMessages(t *testing.T) {
	client := cl.NewClient(":8080")
	server := server.NewServer(":8080", 10*time.Second)
	go server.Start()
	err := client.Connect()
	assert.NoError(t, err)
	err = client.Handshake()
	assert.NoError(t, err)
	checksum, err := client.ReadChecksum()
	assert.NoError(t, err)
	assert.Greater(t, checksum.NumberOfMessages, 0)
	data := make([]int32, 0)
	for i := 0; i < checksum.NumberOfMessages; i++ {
		number, err := client.ReadNumber()
		assert.NoError(t, err)
		data = append(data, number)
	}
	assert.Equal(t, checksum.Checksum, hash.CalculateChecksum(data))
}

func TestReceivingSequenceConnectionDrop(t *testing.T) {
	client := cl.NewClient(":8080")
	server := server.NewServer(":8080", 10*time.Second)
	go server.Start()
	err := client.Connect()
	assert.NoError(t, err)
	err = client.Handshake(cl.HandshakeWithNumberOfMessages(10))
	assert.NoError(t, err)
	checksum, err := client.ReadChecksum()
	assert.NoError(t, err)
	assert.Equal(t, checksum.NumberOfMessages, 10)
	data := make([]int32, 0)
	for i := 0; i < 2; i++ {
		number, err := client.ReadNumber()
		assert.NoError(t, err)
		data = append(data, number)
	}
	err = client.Disconnect()
	assert.NoError(t, err)
	err = client.Connect()
	assert.NoError(t, err)
	err = client.Handshake(cl.HandshakeWithOffset(2))
	assert.NoError(t, err)
	for i := 2; i < 10; i++ {
		number, err := client.ReadNumber()
		assert.NoError(t, err)
		data = append(data, number)
	}
	assert.Equal(t, checksum.Checksum, hash.CalculateChecksum(data))
}
