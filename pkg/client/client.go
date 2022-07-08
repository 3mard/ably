package client

import (
	"encoding/json"
	"fmt"
	"net"

	"3mard.github.com/ably/pkg/message"
	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
)

type Client struct {
	address  string
	conn     net.Conn
	clientId string
	decoder  *json.Decoder
}

type ClientOption func(*Client)

// WithClientId sets the client id.
func WithClientId(id string) ClientOption {
	return func(c *Client) {
		c.clientId = id
	}
}

// NewClient creates a new client.
func NewClient(address string, clientOption ...ClientOption) *Client {
	result := &Client{
		address:  address,
		clientId: uuid.New().String(),
	}

	for _, option := range clientOption {
		option(result)
	}
	return result
}

// Connect connects to the server.
func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}
	c.conn = conn

	c.decoder = json.NewDecoder(c.conn)
	return nil
}

// ConnectWithRetrial retries the connection to the server until it succeeds.
func (c *Client) ConnectWithRetrial(maxRetry uint64) error {
	return backoff.Retry(func() error {
		err := c.Connect()
		if err != nil {
			return err
		}
		return nil
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetry))

}

// Disconnect closes the connection to the server.
func (c *Client) Disconnect() error {
	c.decoder = nil
	return c.conn.Close()
}

type HandshakeOption func(*message.HandshakeMessage)

// HandshakeWithOffset returns a HandshakeOption that sets the offset of the sequence to the given value.
func HandshakeWithOffset(offset int) HandshakeOption {
	return func(handshake *message.HandshakeMessage) {
		handshake.Payload.Offset = offset
		handshake.Payload.Continue = true
	}
}

// HandshakeWithNumberOfMessages returns a HandshakeOption that sets the number of messages to the given value.
func HandshakeWithNumberOfMessages(n int) HandshakeOption {
	return func(handshake *message.HandshakeMessage) {
		handshake.Payload.NumberOfMessages = n
	}
}

// Handshake sends a handshake message to the server.
func (c *Client) Handshake(opts ...HandshakeOption) error {
	handShake := message.HandshakeMessage{
		Type: message.Handshake,
		Payload: message.HandshakePayload{
			UUID: c.clientId,
		},
	}
	for _, opt := range opts {
		opt(&handShake)
	}
	json, err := json.Marshal(handShake)
	if err != nil {
		return err
	}
	c.conn.Write(json)
	return nil
}

// ReadSequence reads a sequence message from the server.
func (c *Client) ReadSequence() (message.SequencePayload, error) {
	var msg message.SequenceMessage
	err := c.decoder.Decode(&msg)
	if err != nil {
		return message.SequencePayload{}, err
	}
	if msg.Type != message.Sequence {
		return message.SequencePayload{}, fmt.Errorf("Expected message of type Sequence, got %d", msg.Type)
	}
	return msg.Payload, nil
}

// ReadChecksum reads a checksum message from the server.
func (c *Client) ReadChecksum() (message.ChecksumPayload, error) {
	var msg message.ChecksumMessage
	err := c.decoder.Decode(&msg)
	if err != nil {
		return message.ChecksumPayload{}, err
	}
	if msg.Type != message.Checksum {
		return message.ChecksumPayload{}, fmt.Errorf("Expected message of type Sequence, got %d", msg.Type)
	}
	return msg.Payload, nil
}

// ReadError reads an error message from the server.
func (c *Client) ReadError() (string, error) {
	var msg message.ErrorMessage
	err := c.decoder.Decode(&msg)
	if err != nil {
		return "", err
	}
	if msg.Type != message.Error {
		return "", fmt.Errorf("Expected message of type Sequence, got %d", msg.Type)
	}
	return msg.Payload, nil
}
