package main

type BaseAgentConfig struct {
	// Client                *LLMClient
	Model string
	// Memory                *Memory
	// SystemPromptGenerator *SystemPromptGenerator
	SystemRole     string
	InputSchema    any
	OutputSchema   any
	Temperature    float32
	MaxTokens      uint64
	ModelAPIParams map[string]string
}

type BaseAgent struct {
	Config *BaseAgentConfig
}
