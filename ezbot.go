package ezbot

import (
	"bufio"
	"fmt"
	"io"
)

type state struct {
	transitions map[string]state
	say         string
	memoryWrite string
	memoryRead  []string
}

var currState, prevState state

var memory map[string]string

func (s *state) buildState(say, memoryWrite string, memoryRead []string) *state {
	newState := state{transitions: make(map[string]state),
		say: say, memoryWrite: memoryWrite, memoryRead: memoryRead}
	if s.transitions != nil {
		s.transitions["default"] = newState
	}
	return &newState
}

// Init initializes the bot.
func Init() {
	memory = make(map[string]string)
	currState = *currState.buildState("Hi. What's your name?", "name", nil)
	currState.buildState("Hello %s! How can I help you?", "des_act", []string{"name"}).
		buildState("Sorry %s, I don't know how to %s. How can I help you?", "des_act", []string{"name", "des_act"})
}

// ConversationLoop runs the bot's conversation loop.
func ConversationLoop(reader io.Reader, writer io.Writer) {
	newReader := bufio.NewReader(reader)
	input := ""
	for {
		fmt.Fprint(writer, act(input))
		input, _ = newReader.ReadString('\n')
		if input == "exit" {
			return
		}
	}
}

func act(input string) string {
	if input != "" {
		memory[prevState.memoryWrite] = input[:len(input)-1]
	}
	memoryOut := retrieveMemory(currState)
	out := fmt.Sprintln()
	out += fmt.Sprintf(currState.say, memoryOut...)

	out += fmt.Sprint(" ")
	prevState = currState
	currState = currState.transitions["default"]
	return out
}

func retrieveMemory(fromState state) []interface{} {
	memoryOut := make([]interface{}, len(currState.memoryRead))
	for i, str := range currState.memoryRead {
		memoryOut[i] = memory[str]
	}
	return memoryOut
}
