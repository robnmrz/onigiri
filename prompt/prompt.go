package prompt

import (
	"fmt"
	"strings"
)

// Struct type for prompt sections
type PromptSection struct {
	Title   string
	Content []string
}

// Option type for SystemPromptGenerator
type SytemPromptGeneratorOption func(*SystemPromptGenerator)

// SystemPromptContextProviderBase interface
type SystemPromptContextProviderBase interface {
	GetInfo() string
	GetTitle() string
}

// SystemPromptGenerator struct
type SystemPromptGenerator struct {
	Background         []string
	Steps              []string
	OutputInstructions []string
	ContextProviders   map[string]SystemPromptContextProviderBase
}

// Constructor for SystemPromptGenerator
// takes in optional function to add background, steps, output instructions and context providers
func NewSystemPromptGenerator(ops ...SytemPromptGeneratorOption) *SystemPromptGenerator {
	spg := &SystemPromptGenerator{
		Background:         []string{},
		Steps:              []string{},
		OutputInstructions: []string{},
		ContextProviders:   map[string]SystemPromptContextProviderBase{},
	}
	for _, opt := range ops {
		opt(spg)
	}
	return spg
}

// Funtions to add a optional background to the system prompt
func WithBackground(background []string) SytemPromptGeneratorOption {
	return func(spg *SystemPromptGenerator) {
		spg.Background = background
	}
}

// Funtions to add optional steps to the system prompt
func WithSteps(steps []string) SytemPromptGeneratorOption {
	return func(spg *SystemPromptGenerator) {
		spg.Steps = steps
	}
}

// Funtions to add optional output instructions to the system prompt
func WithOutputInstructions(outputInstructions []string) SytemPromptGeneratorOption {
	return func(spg *SystemPromptGenerator) {
		spg.OutputInstructions = outputInstructions
	}
}

// Funtions to add optional context providers (runtime input) to the system prompt
func WithContextProviders(contextProviders map[string]SystemPromptContextProviderBase) SytemPromptGeneratorOption {
	return func(spg *SystemPromptGenerator) {
		spg.ContextProviders = contextProviders
	}
}

// GeneratePrompt function to generate the agents system prompt
// based on the available background, steps and output instructions
func (spg *SystemPromptGenerator) GeneratePrompt() string {
	sections := []PromptSection{
		{
			Title:   "IDENTITY and PURPOSE",
			Content: spg.Background,
		},
		{
			Title:   "INTERNAL ASSISTANT STEPS",
			Content: spg.Steps,
		},
		{
			Title:   "OUTPUT INSTRUCTIONS",
			Content: spg.OutputInstructions,
		},
	}

	promptParts := []string{}

	// Add different sections with background and steps and output instructions
	for _, section := range sections {
		if len(section.Content) > 0 {
			promptParts = append(promptParts, fmt.Sprintf("# %s", section.Title))
			for _, content := range section.Content {
				promptParts = append(promptParts, fmt.Sprintf("- %s", content))
			}
			promptParts = append(promptParts, "")
		}
	}

	// Add context providers to prompt
	if len(spg.ContextProviders) > 0 {
		promptParts = append(promptParts, "# EXTRA INFORMATION AND CONTEXT")
		for _, provider := range spg.ContextProviders {
			promptParts = append(promptParts, fmt.Sprintf("# %s", provider.GetTitle()))
			promptParts = append(promptParts, fmt.Sprintf("- %s", provider.GetInfo()))
			promptParts = append(promptParts, "")
		}
	}

	return strings.TrimSpace(strings.Join(promptParts, "\n"))

}
