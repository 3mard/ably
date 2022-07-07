package client

import (
	"encoding/json"
	"fmt"
	"net"

	"3mard.github.com/ably/pkg/message"
	"github.com/google/uuid"
)

type Client struct {
	address  string
	conn     net.Conn
	clientId string
	decoder  *json.Decoder
}

type ClientOption func(*Client)

func WithClientId(id string) ClientOption {
	return func(c *Client) {
		c.clientId = id
	}
}

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

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}
	c.conn = conn

	c.decoder = json.NewDecoder(c.conn)
	return nil
}

func (c *Client) Disconnect() error {
	c.decoder = nil
	return c.conn.Close()
}

type HandshakeOption func(*message.HandshakeMessage)

func HandshakeWithOffset(offset int) HandshakeOption {
	return func(handshake *message.HandshakeMessage) {
		handshake.Payload.Offset = offset
		handshake.Payload.Continue = true
	}
}

func HandshakeWithNumberOfMessages(n int) HandshakeOption {
	return func(handshake *message.HandshakeMessage) {
		handshake.Payload.NumberOfMessages = n
	}
}

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

func (c *Client) ReadNumber() (int32, error) {
	var msg message.SequenceMessage
	err := c.decoder.Decode(&msg)
	if err != nil {
		return 0, err
	}
	if msg.Type != message.Sequence {
		return 0, fmt.Errorf("Expected message of type Sequence, got %d", msg.Type)
	}
	return msg.Payload.Sequence, nil
}

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
