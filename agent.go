package main

import "github.com/robnmrz/onigiri/memory"

type BaseAgentConfig struct {
	// Client                *LLMClient
	Model  string
	Memory *memory.AgentMemory
	// SystemPromptGenerator *SystemPromptGenerator
	SystemRole     string
	InputSchema    any
	OutputSchema   any
	Temperature    float32
	MaxTokens      uint64
	ModelAPIParams map[string]string
}

type BaseAgent struct {
	Config        *BaseAgentConfig
	InitialMemory *memory.AgentMemory
	Memory        *memory.AgentMemory
}

func NewBaseAgent(config *BaseAgentConfig) *BaseAgent {
	return &BaseAgent{
		Config: config,
		Memory: config.Memory.Copy(),
	}
}

// function to
