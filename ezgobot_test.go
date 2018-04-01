package ezgobot

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func test(result, expected string, t *testing.T) {
	if result != expected {
		t.Errorf("Wrong output. Got |%s|, should be |%s|.", result, expected)
	}
}

func TestLoop(t *testing.T) {
	Init()
	w := bytes.NewBufferString("")

	ConversationLoop(strings.NewReader("Bob\nexit"), w)
	test(w.String(), "\nHi. What's your name? \nHello Bob! How can I help you? ", t)
}

func TestAct(t *testing.T) {
	Init()
	test(act(""), "\nHi. What's your name? ", t)
	test(act("Bob\n"), "\nHello Bob! How can I help you? ", t)
	test(act("do nothing\n"), "\nSorry Bob, I don't know how to do nothing. How can I help you? ", t)
	test(act("What is your name?\n"), "\nMy name is Machina. ", t)
}

func TestDetermineTransition(t *testing.T) {
	r1, _ := regexp.Compile("I'm \\d+ years old")
	t1 := transitionMapping{r1, "age_input"}
	r2, _ := regexp.Compile("I'm at \\w")
	t2 := transitionMapping{r2, "loc_input"}
	transitions := []transitionMapping{t1, t2}
	trans := determineTransition("I'm 19 years old", transitions)
	test(trans, "age_input", t)
	trans = determineTransition("I'm at home", transitions)
	test(trans, "loc_input", t)
	trans = determineTransition("ugh", transitions)
	test(trans, "default", t)
}
