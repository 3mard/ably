package message

type MessageType int

const (
	Handshake MessageType = 0
	Sequence  MessageType = 1
	Error     MessageType = 2
	Checksum  MessageType = 3
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

type HandshakePayload struct {
	UUID             string `json:"uuid"`
	Continue         bool   `json:"continue"`
	NumberOfMessages int    `json:"numberOfMessages"`
	Offset           int    `json:"offset"`
}

type HandshakeMessage struct {
	Type    MessageType      `json:"type"`
	Payload HandshakePayload `json:"payload"`
}

type SequencePayload struct {
	Sequence int32 `json:"sequence"`
	Index    int   `json:"index"`
}

type SequenceMessage struct {
	Type    MessageType     `json:"type"`
	Payload SequencePayload `json:"payload"`
}

type ErrorMessage struct {
	Type    MessageType `json:"type"`
	Payload string      `json:"payload"`
}

type ChecksumPayload struct {
	Checksum         string `json:"checksum"`
	NumberOfMessages int    `json:"numberOfMessages"`
}

type ChecksumMessage struct {
	Type    MessageType     `json:"type"`
	Payload ChecksumPayload `json:"payload"`
}
