package main

import (
	"os"

	ezgobot "github.com/ezetter/ezgobot"
)

func buildBot() {
	bootState := ezgobot.Init()
	s1 := bootState.BuildTransitionState("Hi. What's your name?", "default").
		SetMemoryWrite("name").SetID("init")
	s2 := s1.BuildTransitionState("Hello %s! How can I help you?", "default").
		SetID("helpful").
		SetMemoryWrite("des_act").
		SetMemoryRead([]string{"name"})
	s2.BuildTransitionState("Sorry %s, I don't know how to \"%s\". How can I help you?", "default").
		SetID("cantdo").
		SetMemoryRead([]string{"name", "des_act"}).
		AddImmediateTransition(s2)

	s2.BuildTransitionState("My name is Machina. How can I help you?", "ask_name").
		SetID("ask_name").
		AddImmediateTransition(s2)
	s2.AddTransitionMapping("what.*s your name", "ask_name")
}

func main() {
	buildBot()
	ezgobot.ConversationLoop(os.Stdin, os.Stdout)
}
