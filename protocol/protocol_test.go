package protocol

import (
	"bytes"
	"testing"
)

func TestIntToBytes(t *testing.T) {
	x := 256
	actual := intToBytes(x)
	expected := []byte{0, 0, 1, 0}
	checkResult(t, expected, actual)
}

func TestBytesToInt(t *testing.T) {
	b := []byte{0, 0, 0, 1}
	expected := 1
	actual := bytesToInt(b)
	checkResult(t, expected, actual)
}

func TestPack(t *testing.T) {
	msg := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	expected := []byte{'R', 'o', 's', 'e', 0, 0, 0, 6, 'g', 'o', 'l', 'a', 'n', 'g'}
	actual := Pack(msg)
	checkResult(t, expected, actual)
}

func TestUnpack(t *testing.T) {
	msg := []byte{'R', 'o', 's', 'e', 0, 0, 0, 6, 'g', 'o', 'l', 'a', 'n', 'g'}
	expected := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	r := make(chan []byte, 16)
	Unpack(msg, r)
	actual := <-r
	checkResult(t, expected, actual)
}

func checkResult(t *testing.T, expected, actual interface{}) {
	var result bool
	switch expected.(type) {
	case []byte:
		e := expected.([]byte)
		a := actual.([]byte)
		result = bytes.Equal(e, a)
	default:
		result = expected == actual
	}
	if !result {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}
