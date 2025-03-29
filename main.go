package main

import (
	"fmt"

	"github.com/robnmrz/onigiri/memory"
)

type AgentMessageContent struct {
	FileName string `json:"file_name"`
	Message  string `json:"message"`
}

func main() {
	messageTest := AgentMessageContent{
		FileName: "test.txt",
		Message:  "Hello world",
	}
	memory := memory.NewAgentMemory(memory.WithMaxMessages(5))
	memory.InitializeTurn()
	memory.AddMessage("user", messageTest)

	jsonString, err := memory.ToJson()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonString)
}
