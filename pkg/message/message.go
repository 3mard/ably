package message

type MessageType int

const (
	Handshake MessageType = iota
	Sequence  MessageType = iota
	Error     MessageType = iota
	Checksum  MessageType = iota
)

type Message struct {
	Type    MessageType
	Payload interface{}
}

type HandshakePayload struct {
	UUID             string
	Continue         bool
	NumberOfMessages int
	Offset           int
}

type HandshakeMessage struct {
	Type    MessageType
	Payload HandshakePayload
}

type SequencePayload struct {
	Sequence int32
	Index    int
}

type SequenceMessage struct {
	Type    MessageType
	Payload SequencePayload
}

type ErrorMessage struct {
	Type    MessageType
	Payload string
}

type ChecksumPayload struct {
	Checksum         string
	NumberOfMessages int
}

type ChecksumMessage struct {
	Type    MessageType
	Payload ChecksumPayload
}
