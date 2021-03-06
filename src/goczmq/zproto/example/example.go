// Package example is 100% generated. If you edit this file,
// you will lose your changes at the next build cycle.
// DO NOT MAKE ANY CHANGES YOU WISH TO KEEP.
//
// The correct places for commits are:
//  - The XML model used for this code generation: zproto_example_goczmq.xml
//  - The code generation script that built this file: zproto_codec_go
package example

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/zeromq/goczmq"
)

const (
	Signature uint16 = 0xAAA0 | 0
	Version          = 1
)

const (
	LogId        uint8 = 1
	StructuresId uint8 = 2
	BinaryId     uint8 = 3
	TypesId      uint8 = 4
)

// Transit is a codec interface
type Transit interface {
	Marshal() ([]byte, error)
	Unmarshal(...[]byte) error
	String() string
	Send(*goczmq.Sock) error
	SetRoutingId([]byte)
	RoutingId() []byte
	SetSequence(uint16)
	Sequence() uint16
}

// Unmarshal unmarshals data from raw frames.
func Unmarshal(frames ...[]byte) (t Transit, err error) {
	if frames == nil {
		return nil, errors.New("can't unmarshal an empty message")
	}
	var buffer *bytes.Buffer

	// Check the signature
	var signature uint16
	buffer = bytes.NewBuffer(frames[0])
	binary.Read(buffer, binary.BigEndian, &signature)
	if signature != Signature {
		// Invalid signature
		return nil, fmt.Errorf("invalid signature %X != %X", Signature, signature)
	}

	// Get message id and parse per message type
	var id uint8
	binary.Read(buffer, binary.BigEndian, &id)

	switch id {
	case LogId:
		t = NewLog()
	case StructuresId:
		t = NewStructures()
	case BinaryId:
		t = NewBinary()
	case TypesId:
		t = NewTypes()
	}
	err = t.Unmarshal(frames...)

	return t, err
}

// Recv receives marshaled data from a 0mq socket.
func Recv(sock *goczmq.Sock) (t Transit, err error) {
	return recv(sock, 0)
}

// RecvNoWait receives marshaled data from 0mq socket. It won't wait for input.
func RecvNoWait(sock *goczmq.Sock) (t Transit, err error) {
	return recv(sock, goczmq.DONTWAIT)
}

// recv receives marshaled data from 0mq socket.
func recv(sock *goczmq.Sock, flag goczmq.Flag) (t Transit, err error) {
	var frames [][]byte

	if flag == goczmq.DONTWAIT {
		frames, err = sock.RecvMessageNoWait()
	} else {
		frames, err = sock.RecvMessage()
	}

	if err != nil {
		return nil, err
	}

	sType := sock.GetType()
	if err != nil {
		return nil, err
	}

	var routingId []byte
	// If message came from a router socket, first frame is routingId
	if sType == goczmq.ROUTER {
		if len(frames) <= 1 {
			return nil, errors.New("no routingId")
		}
		routingId = frames[0]
		frames = frames[1:]
	}

	t, err = Unmarshal(frames...)
	if err != nil {
		return nil, err
	}

	if sType == goczmq.ROUTER {
		t.SetRoutingId(routingId)
	}
	return t, err
}

// Clone clones a message.
func Clone(t Transit) Transit {

	switch msg := t.(type) {
	case *Log:
		cloned := NewLog()
		routingId := make([]byte, len(msg.RoutingId()))
		copy(routingId, msg.RoutingId())
		cloned.SetRoutingId(routingId)
		cloned.sequence = msg.sequence
		cloned.Version = msg.Version
		cloned.Level = msg.Level
		cloned.Event = msg.Event
		cloned.Node = msg.Node
		cloned.Peer = msg.Peer
		cloned.Time = msg.Time
		cloned.Host = msg.Host
		cloned.Data = msg.Data
		return cloned

	case *Structures:
		cloned := NewStructures()
		routingId := make([]byte, len(msg.RoutingId()))
		copy(routingId, msg.RoutingId())
		cloned.SetRoutingId(routingId)
		cloned.sequence = msg.sequence
		for idx, str := range msg.Aliases {
			cloned.Aliases[idx] = str
		}
		for key, val := range msg.Headers {
			cloned.Headers[key] = val
		}
		return cloned

	case *Binary:
		cloned := NewBinary()
		routingId := make([]byte, len(msg.RoutingId()))
		copy(routingId, msg.RoutingId())
		cloned.SetRoutingId(routingId)
		cloned.sequence = msg.sequence
		cloned.Flags = msg.Flags
		cloned.PublicKey = append(cloned.PublicKey, msg.PublicKey...)
		cloned.Identifier = append(cloned.Identifier, msg.Identifier...)
		cloned.Address = append(cloned.Address, msg.Address...)
		cloned.Content = append(cloned.Content, msg.Content...)
		return cloned

	case *Types:
		cloned := NewTypes()
		routingId := make([]byte, len(msg.RoutingId()))
		copy(routingId, msg.RoutingId())
		cloned.SetRoutingId(routingId)
		cloned.sequence = msg.sequence
		cloned.ClientForename = msg.ClientForename
		cloned.ClientSurname = msg.ClientSurname
		cloned.ClientMobile = msg.ClientMobile
		cloned.ClientEmail = msg.ClientEmail
		cloned.SupplierForename = msg.SupplierForename
		cloned.SupplierSurname = msg.SupplierSurname
		cloned.SupplierMobile = msg.SupplierMobile
		cloned.SupplierEmail = msg.SupplierEmail
		return cloned
	}

	return nil
}

// putString marshals a string into the buffer.
func putString(buffer *bytes.Buffer, str string) {
	size := len(str)
	binary.Write(buffer, binary.BigEndian, byte(size))
	binary.Write(buffer, binary.BigEndian, []byte(str[0:size]))
}

// getString unmarshals a string from the buffer.
func getString(buffer *bytes.Buffer) string {
	var size byte
	binary.Read(buffer, binary.BigEndian, &size)
	str := make([]byte, size)
	binary.Read(buffer, binary.BigEndian, &str)
	return string(str)
}

// putLongString marshals a string into the buffer.
func putLongString(buffer *bytes.Buffer, str string) {
	size := len(str)
	binary.Write(buffer, binary.BigEndian, uint32(size))
	binary.Write(buffer, binary.BigEndian, []byte(str[0:size]))
}

// getLongString unmarshals a string from the buffer.
func getLongString(buffer *bytes.Buffer) string {
	var size uint32
	binary.Read(buffer, binary.BigEndian, &size)
	str := make([]byte, size)
	binary.Read(buffer, binary.BigEndian, &str)
	return string(str)
}

// putBytes marshals []byte into the buffer.
func putBytes(buffer *bytes.Buffer, data []byte) {
	size := uint64(len(data))
	binary.Write(buffer, binary.BigEndian, size)
	binary.Write(buffer, binary.BigEndian, data)
}

// getBytes unmarshals []byte from the buffer.
func getBytes(buffer *bytes.Buffer) []byte {
	var size uint64
	binary.Read(buffer, binary.BigEndian, &size)
	data := make([]byte, size)
	binary.Read(buffer, binary.BigEndian, &data)
	return data
}
