package agent

import (
	"context"
	"errors"
	"fmt" // Using log for deprecation warnings, similar to Python's warnings
	"reflect"

	"github.com/robnmrz/onigiri/memory"
	"github.com/robnmrz/onigiri/prompt"
)

var (
	// DefaultAgentInputSchema represents the default type for agent input.
	DefaultAgentInputSchema reflect.Type = reflect.TypeOf("") // Example: expecting string input
	// DefaultAgentOutputSchema represents the default type for agent output.
	DefaultAgentOutputSchema reflect.Type = reflect.TypeOf("") // Example: expecting string output
)

// LLMClient defines the interface for interacting with a Language Model.
type LLMClient interface {
	CreateCompletion(messages []memory.Message, responseSchema reflect.Type, model string, modelApiParameters map[string]any) (CompletionResponse, error)
}

type CompletionResponse struct {
	Prompt string
}

// AgentConfig holds the configuration for BaseAgent, applied via options.
type AgentConfig struct {
	inputSchema           reflect.Type
	outputSchema          reflect.Type
	client                LLMClient
	model                 string
	memory                *memory.AgentMemory
	systemPromptGenerator *prompt.SystemPromptGenerator
	systemRole            string
	modelApiParameters    map[string]any
}

// AgentOption defines the functional option type.
type AgentOption func(*AgentConfig) error

// BaseAgent is the core implementation for chat agents.
type BaseAgent struct {
	client                LLMClient
	model                 string
	memory                *memory.AgentMemory
	initialMemory         *memory.AgentMemory
	systemPromptGenerator *prompt.SystemPromptGenerator
	systemRole            string
	modelApiParameters    map[string]any
	// inputSchema           reflect.Type
	outputSchema     reflect.Type
	currentUserInput any
}

// WithInputSchema sets the expected input type for the agent.
func WithInputSchema(schemaType reflect.Type) AgentOption {
	return func(cfg *AgentConfig) error {
		if schemaType == nil {
			return errors.New("input schema type cannot be nil")
		}
		// Basic kind check (can be enhanced)
		if schemaType.Kind() == reflect.Invalid || schemaType.Kind() == reflect.Func {
			return fmt.Errorf("invalid kind for input schema: %v", schemaType.Kind())
		}
		cfg.inputSchema = schemaType
		return nil
	}
}

// WithOutputSchema sets the expected output type for the agent.
func WithOutputSchema(schemaType reflect.Type) AgentOption {
	return func(cfg *AgentConfig) error {
		if schemaType == nil {
			return errors.New("output schema type cannot be nil")
		}
		// Basic kind check (can be enhanced)
		if schemaType.Kind() == reflect.Invalid || schemaType.Kind() == reflect.Func {
			return fmt.Errorf("invalid kind for output schema: %v", schemaType.Kind())
		}
		cfg.outputSchema = schemaType
		return nil
	}
}

// WithClient sets the LLM client dependency. This is mandatory.
func WithClient(client LLMClient) AgentOption {
	return func(cfg *AgentConfig) error {
		if client == nil {
			return errors.New("LLM client cannot be nil")
		}
		cfg.client = client
		return nil
	}
}

// WithModel sets the LLM model name. This is mandatory.
func WithModel(model string) AgentOption {
	return func(cfg *AgentConfig) error {
		if model == "" {
			return errors.New("model name cannot be empty")
		}
		cfg.model = model
		return nil
	}
}

// WithMemory sets the agent's memory component.
// If not provided, a default implementation might be used in NewBaseAgent.
func WithMemory(memory *memory.AgentMemory) AgentOption {
	return func(cfg *AgentConfig) error {
		cfg.memory = memory
		return nil
	}
}

// WithSystemPromptGenerator sets the system prompt generator.
// If not provided, a default implementation might be used in NewBaseAgent.
func WithSystemPromptGenerator(spg *prompt.SystemPromptGenerator) AgentOption {
	return func(cfg *AgentConfig) error {
		cfg.systemPromptGenerator = spg
		return nil
	}
}

// WithSystemRole sets the role name for the system prompt message (e.g., "system").
func WithSystemRole(role string) AgentOption {
	return func(cfg *AgentConfig) error {
		// Allow empty role to disable system prompt message explicitly
		cfg.systemRole = role
		return nil
	}
}

// WithModelParameter adds a key-value pair to the additional API parameters.
func WithModelParameter(key string, value any) AgentOption {
	return func(cfg *AgentConfig) error {
		if key == "" {
			return errors.New("model API parameter key cannot be empty")
		}
		if cfg.modelApiParameters == nil {
			cfg.modelApiParameters = make(map[string]any)
		}
		cfg.modelApiParameters[key] = value
		return nil
	}
}

// NewBaseAgent creates a new BaseAgent instance with the provided options.
func NewBaseAgent(opts ...AgentOption) (*BaseAgent, error) {
	// Initialize config with defaults
	cfg := &AgentConfig{
		client:                nil,
		inputSchema:           DefaultAgentInputSchema, // Use package-level defaults
		outputSchema:          DefaultAgentOutputSchema,
		modelApiParameters:    make(map[string]any),
		systemRole:            "system",
		memory:                nil,
		systemPromptGenerator: nil,
	}

	// Apply all options
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, fmt.Errorf("failed to apply agent option: %w", err)
		}
	}

	// Create the agent
	agent := &BaseAgent{
		client:                cfg.client,
		model:                 cfg.model,
		memory:                cfg.memory,
		systemPromptGenerator: cfg.systemPromptGenerator,
		systemRole:            cfg.systemRole,
		modelApiParameters:    cfg.modelApiParameters,
	}

	// Store the initial memory state for resets
	agent.initialMemory = agent.memory.Copy()

	return agent, nil
}

// ResetMemory resets the agent's memory to its initial state.
func (a *BaseAgent) ResetMemory() {
	a.memory = a.initialMemory.Copy()
}

func (a *BaseAgent) GetResponse() (CompletionResponse, error) {
	var messages []memory.Message
	responseModel := a.outputSchema

	// Omit system prompt if role is empty
	if a.systemRole == "" {
		messages = []memory.Message{}
	} else {
		messages = []memory.Message{
			{
				Role: a.systemRole,
				Content: memory.MessageContent{
					TypeName: "string",
					Content:  a.systemPromptGenerator.GeneratePrompt(),
				},
			},
		}
	}

	// Add messages from memory
	messages = append(messages, a.memory.History...)

	response, err := a.client.CreateCompletion(messages, responseModel, a.model, a.modelApiParameters)
	if err != nil {
		return CompletionResponse{}, err
	}
	return response, nil
}

func (a *BaseAgent) Run(ctx context.Context, userInput any) (CompletionResponse, error) {

	// Init a new turn when user gives input
	if userInput != nil {
		a.memory.InitializeTurn()
		a.memory.AddMessage("user", userInput)
		a.currentUserInput = userInput
	}

	response, err := a.GetResponse()
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("LLM completion failed: %w", err) // Error already includes context from GetResponse
	}

	// Add assistant response to memory
	a.memory.AddMessage("assistant", response)
	return response, nil
}

// GetContextProvider retrieves a context provider by name from the SystemPromptGenerator.
func (a *BaseAgent) GetContextProvider(providerName string) (prompt.SystemPromptContextProviderBase, error) {
	if a.systemPromptGenerator == nil {
		return nil, errors.New("system prompt generator is not configured")
	}

	provider, ok := a.systemPromptGenerator.ContextProviders[providerName]
	if !ok {
		return nil, fmt.Errorf("context provider '%s' not found", providerName)
	}
	return provider, nil
}

// RegisterContextProvider registers a new context provider with the SystemPromptGenerator.
func (a *BaseAgent) RegisterContextProvider(providerName string, provider prompt.SystemPromptContextProviderBase) error {
	if a.systemPromptGenerator == nil {
		return errors.New("system prompt generator is not configured")
	}

	if providerName == "" {
		return errors.New("provider name cannot be empty")
	}
	a.systemPromptGenerator.ContextProviders[providerName] = provider
	return nil
}

// UnregisterContextProvider removes a context provider from the SystemPromptGenerator.
func (a *BaseAgent) UnregisterContextProvider(providerName string) error {
	if a.systemPromptGenerator == nil {
		return errors.New("system prompt generator is not configured")
	}
	if providerName == "" {
		return errors.New("provider name cannot be empty")
	}

	delete(a.systemPromptGenerator.ContextProviders, providerName)
	return nil
}
