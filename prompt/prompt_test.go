package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock context provider
type MockContextProvider struct {
	mock.Mock
}

func (m *MockContextProvider) GetInfo() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockContextProvider) GetTitle() string {
	args := m.Called()
	return args.String(0)
}

func TestNewSystemPromptGenerator_Default(t *testing.T) {
	spg := NewSystemPromptGenerator()
	assert.NotNil(t, spg)
	assert.Empty(t, spg.Background)
	assert.Empty(t, spg.Steps)
	assert.Empty(t, spg.OutputInstructions)
	assert.Empty(t, spg.ContextProviders)
}

func TestWithBackground(t *testing.T) {
	bg := []string{"I am an assistant", "My job is to help users"}
	spg := NewSystemPromptGenerator(WithBackground(bg))

	assert.Equal(t, bg, spg.Background)
}

func TestWithSteps(t *testing.T) {
	steps := []string{"Step 1", "Step 2"}
	spg := NewSystemPromptGenerator(WithSteps(steps))

	assert.Equal(t, steps, spg.Steps)
}

func TestWithOutputInstructions(t *testing.T) {
	instructions := []string{"Format the output as JSON"}
	spg := NewSystemPromptGenerator(WithOutputInstructions(instructions))

	assert.Equal(t, instructions, spg.OutputInstructions)
}

func TestWithContextProviders(t *testing.T) {
	mockProvider := new(MockContextProvider)
	contextMap := map[string]SystemPromptContextProviderBase{
		"provider1": mockProvider,
	}
	spg := NewSystemPromptGenerator(WithContextProviders(contextMap))

	assert.Equal(t, contextMap, spg.ContextProviders)
}

func TestGeneratePrompt_AllSections(t *testing.T) {
	mockProvider := new(MockContextProvider)
	mockProvider.On("GetTitle").Return("User Info")
	mockProvider.On("GetInfo").Return("This is some extra context")

	spg := NewSystemPromptGenerator(
		WithBackground([]string{"I am a helpful assistant."}),
		WithSteps([]string{"Greet the user", "Answer the question"}),
		WithOutputInstructions([]string{"Respond in markdown format."}),
		WithContextProviders(map[string]SystemPromptContextProviderBase{
			"mock": mockProvider,
		}),
	)

	prompt := spg.GeneratePrompt()

	assert.Contains(t, prompt, "# IDENTITY and PURPOSE")
	assert.Contains(t, prompt, "- I am a helpful assistant.")
	assert.Contains(t, prompt, "# INTERNAL ASSISTANT STEPS")
	assert.Contains(t, prompt, "- Greet the user")
	assert.Contains(t, prompt, "# OUTPUT INSTRUCTIONS")
	assert.Contains(t, prompt, "- Respond in markdown format.")
	assert.Contains(t, prompt, "# EXTRA INFORMATION AND CONTEXT")
	assert.Contains(t, prompt, "# User Info")
	assert.Contains(t, prompt, "- This is some extra context")
}
