package ezbot

import (
	"bytes"
	"strings"
	"testing"
)

func test(expected, result string, t *testing.T) {
	if result != expected {
		t.Errorf("Wrong output. Got |%s|, should be |%s|.", result, expected)
	}
}

func TestLoop(t *testing.T) {
	Init()
	w := bytes.NewBufferString("")

	ConversationLoop(strings.NewReader("Bob\nexit"), w)
	test("\nHi. What's your name? \nHello Bob! How can I help you? ", w.String(), t)
}

func TestAct(t *testing.T) {
	Init()
	test(act(""), "\nHi. What's your name? ", t)
	test(act("Bob\n"), "\nHello Bob! How can I help you? ", t)
	test(act("do nothing\n"), "\nSorry Bob, I don't know how to do nothing. How can I help you? ", t)
}
