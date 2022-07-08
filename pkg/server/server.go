package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"3mard.github.com/ably/pkg/hash"
	"3mard.github.com/ably/pkg/message"
)

const (
	MaxNumberOfMessages = 0xffff
)

type Server struct {
	address    string
	clientRepo ClientRepository
}

// NewServer creates a new server
func NewServer(address string, ttl time.Duration) *Server {
	return &Server{address: address,
		clientRepo: NewInMemoryClientRepository(ttl)}
}

// Start starts the server
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
		go s.handelConnection(conn)
	}
}

func (s *Server) handelConnection(conn net.Conn) {
	msg, err := s.handleHandshake(conn)
	if err != nil {
		log.Println("error handling handshake: ", err)
	}
	id := msg.Payload.UUID
	var numbersToSend []int32
	offset := 0
	if msg.Payload.NumberOfMessages > MaxNumberOfMessages {
		errMessage := message.ErrorMessage{
			Type:    message.Error,
			Payload: fmt.Sprintf("Number of messages is too big: %d", msg.Payload.NumberOfMessages),
		}
		s.sendMessage(conn, &errMessage)
		return
	}
	if msg.Payload.NumberOfMessages == 0 {
		msg.Payload.NumberOfMessages = rand.Int() % MaxNumberOfMessages
	}

	if !msg.Payload.Continue {
		numbersToSend = s.generateNumbersToSend(msg.Payload.NumberOfMessages)
		err := s.clientRepo.SetClientMessages(id, numbersToSend)
		if err != nil {
			fmt.Println("Error saving messages to db", err)
		}
		checksum := hash.CalculateChecksum(numbersToSend)
		checksumMessage := message.ChecksumMessage{
			Type: message.Checksum,
			Payload: message.ChecksumPayload{
				Checksum:         checksum,
				NumberOfMessages: len(numbersToSend),
			},
		}
		err = s.sendMessage(conn, checksumMessage)
	} else {
		numbersToSend, err = s.clientRepo.GetClientMessages(id)
		if err != nil {
			msg := message.ErrorMessage{
				Type:    message.Error,
				Payload: "Can't continue progress, please start over",
			}
			s.sendMessage(conn, msg)
			return
		}
		if msg.Payload.Offset > len(numbersToSend) {
			msg := message.ErrorMessage{
				Type:    message.Error,
				Payload: "Offset is out of range, Please start over",
			}
			s.sendMessage(conn, msg)
			return
		}
		offset = msg.Payload.Offset
	}

	err = s.sendMessages(offset, numbersToSend, conn)
	if err != nil {
		fmt.Println("Error sending messages:", err)
	}

}

func (s *Server) handleHandshake(conn net.Conn) (message.HandshakeMessage, error) {
	buf := make([]byte, 1024)
	r := bufio.NewReader(conn)
	msgSize, err := r.Read(buf)
	if err != nil {
		return message.HandshakeMessage{}, err
	}

	var handShakeMessage message.HandshakeMessage
	err = json.Unmarshal(buf[:msgSize], &handShakeMessage)
	if err != nil {
		log.Println("Error: couldn't unmarshal", err)
		return message.HandshakeMessage{}, err
	}
	if handShakeMessage.Type != message.Handshake {
		return message.HandshakeMessage{}, fmt.Errorf("Expected handshake message")
	}

	log.Println("shook hands")
	return handShakeMessage, nil
}

func (s *Server) generateNumbersToSend(n int) []int32 {
	result := make([]int32, n)
	for i := 0; i < n; i++ {
		result[i] = int32(rand.Uint32())
	}
	return result
}

func (s *Server) sendMessages(offset int, msgs []int32, conn net.Conn) error {
	for i := offset; i < len(msgs); i++ {
		seqMessage := message.SequenceMessage{
			Type: message.Sequence,
			Payload: message.SequencePayload{
				Sequence: msgs[i],
				Index:    i,
			},
		}
		s.sendMessage(conn, seqMessage)
	}
	return nil
}

func (s *Server) sendMessage(con net.Conn, data interface{}) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = con.Write(json)
	if err != nil {
		return err
	}
	return nil
}
