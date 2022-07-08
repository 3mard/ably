
## Message:
This protocol is json based, every message exchanged between the server and the client is encoded in the json format
The message two properties :
- Type: specifies  the type of the message 
- payload: specifies the payload of the message

The protocol supports the following messages:
- Handshake(type:0) A message sent from the client to server to establish connection which, specifies the UUID of the client, How many messages does the client wish for (max 0xffff), whether the client wishes to resume the flow from a cretin offset, and what the offset might be
- Checksum(type:3) A message sent from the server to the client, which tells who many messages should the client expect and, what is the combined hash for the messages
- Sequence(type:1) A message sent from the server to the client, which specifies the offset of this sequence item,  and the value of it
- Error(type:2) A message sent from the server to client, which indelicate that something has gone wrong

## The communication flow 

- The client send a handshake message to the server
- The server acknowledge the handshake with a checksum message
- The server starts sending sequence messages to the client
- The server closes the connection once the messages has been sent
- The server can send an error message at any point of time

## How to build
To build the client run:
```sh
make build-client
```
To build the server run:
```sh
make build-server
```

## How to run
To run the client
```sh
./client -address=<address:port>
```
eg:
```sh
./client -address=localhost:8080
```

To run the server
```sh
./server -address=<address:port>
```
eg:
```sh
./server -address=localhost:8080
```
To run the tests:
```sh
make tests
```

## ToDos
- Improve how the client handles the errors, the client almost never check if the message is an error message or the expected message
- Make sure no more 1 connection with the same UUID is established.
- Add unit tests
- The client assume that the messages will arrive in the right order, however the protocol makes no such assumptions 
- The server should implement a graceful shutdown 
- Use a proper logger
