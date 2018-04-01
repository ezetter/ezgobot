package ezbot

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

type transitionMapping struct {
	inputMatch *regexp.Regexp
	transition string
}

// State defines a node in the state machine representing a conversation.
type State struct {
	transitions             map[string]*State
	transitionInputMappings []transitionMapping
	say                     string
	memoryWrite             string
	memoryRead              []string
}

var currState, prevState State

var memory map[string]string

// BuildTransitionState builds a state as a transition from the reciever
func (s *State) BuildTransitionState(say string, transitionName string) *State {
	newState := &State{transitions: make(map[string]*State), say: say}
	if s.transitions != nil {
		s.transitions["default"] = newState
	}
	return newState
}

// AddTransitionMapping adds a mapping from inputs to transitions.
func (s *State) AddTransitionMapping(matchRegExp, transition string) *State {
	inputMatch, _ := regexp.Compile(matchRegExp)
	s.transitionInputMappings = append(s.transitionInputMappings, transitionMapping{inputMatch, transition})
	return s
}

// SetMemoryWrite sets the memory that will be written after the state runs.
func (s *State) SetMemoryWrite(memoryWrite string) *State {
	s.memoryWrite = memoryWrite
	return s
}

// SetMemoryRead sets the memory elements that the state needs to read.
func (s *State) SetMemoryRead(memoryRead []string) *State {
	s.memoryRead = memoryRead
	return s
}

// AddTransition adds a new state transition.
func (s *State) AddTransition(transitionName string, destination *State) *State {
	s.transitions[transitionName] = destination
	return s
}

// Init initializes the bot.
func Init() {
	memory = make(map[string]string)
	currState = *currState.BuildTransitionState("Hi. What's your name?", "default").
		SetMemoryWrite("name")
	s2 := currState.BuildTransitionState("Hello %s! How can I help you?", "default").
		SetMemoryWrite("des_act").
		SetMemoryRead([]string{"name"}).
		BuildTransitionState("Sorry %s, I don't know how to %s. How can I help you?", "default").
		SetMemoryWrite("des_act").
		SetMemoryRead([]string{"name", "des_act"})
	s2.AddTransition("default", s2)
	s2.BuildTransitionState("My name is Machina.", "ask_name").AddTransition("default", s2)
	s2.AddTransitionMapping("What is your name", "ask_name")
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

func determineTransition(input string, mappings []transitionMapping) string {
	for _, mapping := range mappings {
		if mapping.inputMatch.MatchString(input) {
			return mapping.transition
		}
	}
	return "default"
}

func act(input string) string {
	if input != "" {
		memory[prevState.memoryWrite] = input[:len(input)-1]
	}
	remembered := retrieveMemory(currState)
	out := fmt.Sprintln()
	out += fmt.Sprintf(currState.say, remembered...)

	out += fmt.Sprint(" ")
	prevState = currState
	currState = *currState.transitions["default"]
	return out
}

func retrieveMemory(fromState State) []interface{} {
	memoryOut := make([]interface{}, len(currState.memoryRead))
	for i, str := range currState.memoryRead {
		memoryOut[i] = memory[str]
	}
	return memoryOut
}
