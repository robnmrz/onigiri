package memory

import (
	"encoding/json"
	"fmt"

	"slices"

	"github.com/google/uuid"
	"github.com/robnmrz/onigiri/utils"
)

type MemoryOption func(*AgentMemory)

// Image is a struct that hold the values and keys of images
type Image struct {
	Value string `json:"url"`
	Key   string `json:"key"`
}

type MessageContent struct {
	TypeName string `json:"type_name"`
	Content  any    `json:"content"`
}

// Message is a struct that holds the role and content of a message
type Message struct {
	Role    string         `json:"role"`
	Content MessageContent `json:"content"`
	TurnId  string         `json:"turn_id"`
}

// TODO: Maybe implementing Messages a map of turnId and Message
// AgentMemory is a struct that holds the memory of the agent
// in form of the chat history, the current turn id and the max messages
type AgentMemory struct {
	History       []Message `json:"history"`
	MaxMessages   int       `json:"max_messages"`
	CurrentTurnId string    `json:"current_turn_id"`
}

// Constructor for a new AgentMemory struct with options
func NewAgentMemory(ops ...MemoryOption) *AgentMemory {
	am := &AgentMemory{
		History:     []Message{},
		MaxMessages: -1,
	}

	for _, op := range ops {
		op(am)
	}

	return am
}

// Functional option to set MaxMessages, otherwise default to -1.
// Defines the maximum number of messages in the history
func WithMaxMessages(maxMessages int) MemoryOption {
	return func(am *AgentMemory) {
		am.MaxMessages = maxMessages
	}
}

// Initialize a new turn
func (am *AgentMemory) InitializeTurn() {
	am.CurrentTurnId = uuid.New().String()
}

// Add a new message to the history
func (am *AgentMemory) AddMessage(role string, content any) {
	am.History = append(am.History, Message{
		Role:    role,
		Content: MessageContent{TypeName: utils.GetTypeName(content), Content: content},
		TurnId:  am.CurrentTurnId,
	})

	// Handle overflow
	am.manageOverflow()
}

// if MaxMessages is set, remove old messages if
// there are more than MaxMessages
func (am *AgentMemory) manageOverflow() {
	if am.MaxMessages != -1 {
		for i := range len(am.History) - am.MaxMessages {
			am.History = am.History[i+1:]
		}
	}
}

// Get all the messages in the history
// handling both regular and multi modal content
// func (am *AgentMemory) GetMessages() ([]Message, error) {
// 	for _, msg := range am.History {
// 		images := []Image{}
// 		contentJson, err := msg.Content.ToJson(
// 		if err != nil {
// 			return []Message{}, err
// 		}
// 	}
// }

// Copy the memory to a new struct
func (am *AgentMemory) Copy() *AgentMemory {
	return &AgentMemory{
		History:       am.History,
		MaxMessages:   am.MaxMessages,
		CurrentTurnId: am.CurrentTurnId,
	}
}

// Getter for retrieving the current turn id
func (am *AgentMemory) GetTurnId() string {
	return am.CurrentTurnId
}

// Function to delete messages by turn id. If message for turn Id
// is not found, return error
func (am *AgentMemory) DeleteMessagesByTurnId(turnId string) error {
	initialLength := len(am.History)
	for i, msg := range am.History {
		if msg.TurnId == turnId {
			am.History = slices.Delete(am.History, i, i+1)
		}
	}
	if len(am.History) == initialLength {
		return fmt.Errorf("message with turn id %s not found", turnId)
	}
	return nil
}

// Function to get the number of messages in the history
func (am *AgentMemory) GetMessageCount() int {
	return len(am.History)
}

// Function to serialize the memory to json.
// Returns the json string or returns an error
func (am *AgentMemory) ToJson() (string, error) {
	jsonBytes, err := json.MarshalIndent(am, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// Function to deserialize the memory from json.
// Populated the memory with the json data or returns an error
func (am *AgentMemory) FromJson(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), am)
}
