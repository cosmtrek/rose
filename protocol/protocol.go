package protocol

import (
	"bytes"
	"encoding/binary"
)

const (
	PackageHeader        = "Rose"
	PackageHeaderLength  = 4
	PackageMinDataLength = 4
)

func Pack(message []byte) []byte {
	return append(append([]byte(PackageHeader), intToBytes(len(message))...), message...)
}

func Unpack(message []byte, reader chan<- []byte) []byte {
	length := len(message)

	var i int
	for i = 0; i < length; i++ {
		if length < i+PackageHeaderLength+PackageMinDataLength {
			break
		}
		if string(message[i:i+PackageHeaderLength]) == PackageHeader {
			messageLen := bytesToInt(message[i+PackageHeaderLength : i+PackageHeaderLength+PackageMinDataLength])
			if length < i+PackageHeaderLength+PackageMinDataLength+messageLen {
				break
			}
			data := message[i+PackageHeaderLength+PackageMinDataLength : i+PackageHeaderLength+PackageMinDataLength+messageLen]
			reader <- data
			i += PackageHeaderLength + PackageMinDataLength + messageLen - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return message[i:]
}

func intToBytes(n int) []byte {
	x := int32(n)
	b := bytes.NewBuffer([]byte{})
	binary.Write(b, binary.BigEndian, x)
	return b.Bytes()
}

func bytesToInt(b []byte) int {
	nb := bytes.NewBuffer(b)
	var x int32
	binary.Read(nb, binary.BigEndian, &x)
	return int(x)
}
