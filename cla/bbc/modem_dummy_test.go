package bbc

import (
	"bytes"
	"testing"
)

func TestDummyModemFragments(t *testing.T) {
	var hub = newDummyHub()
	var modems [10]Modem

	for i := 0; i < len(modems); i++ {
		modems[i] = newDummyModem(16, hub)
	}

	msg := []byte("hello world")
	_ = modems[0].Send(NewFragment(0, 0, false, false, msg))
	for i := 0; i < len(modems); i++ {
		if f, _ := modems[i].Receive(); bytes.Compare(f.Payload, msg) != 0 {
			t.Fatalf("Wrong payload: expected %x, got %x", msg, f.Payload)
		}
	}
}
