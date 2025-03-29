package memory

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DummyContent struct {
	Text string `json:"text"`
}

func TestNewAgentMemory_Defaults(t *testing.T) {
	am := NewAgentMemory()

	assert.NotNil(t, am)
	assert.Equal(t, 0, len(am.History))
	assert.Equal(t, -1, am.MaxMessages)
}

func TestWithMaxMessages(t *testing.T) {
	am := NewAgentMemory(WithMaxMessages(3))

	assert.Equal(t, 3, am.MaxMessages)
}

func TestInitializeTurn(t *testing.T) {
	am := NewAgentMemory()
	am.InitializeTurn()

	assert.NotEmpty(t, am.CurrentTurnId)
}

func TestAddMessage(t *testing.T) {
	am := NewAgentMemory()
	am.InitializeTurn()
	am.AddMessage("user", DummyContent{Text: "Hello"})

	assert.Equal(t, 1, len(am.History))
	assert.Equal(t, "user", am.History[0].Role)
	assert.Equal(t, am.CurrentTurnId, am.History[0].TurnId)
	assert.Equal(t, "DummyContent", am.History[0].Content.TypeName)
}

func TestMessageOverflow(t *testing.T) {
	am := NewAgentMemory(WithMaxMessages(2))
	am.InitializeTurn()
	am.AddMessage("user", DummyContent{Text: "msg1"})
	am.AddMessage("user", DummyContent{Text: "msg2"})
	am.AddMessage("user", DummyContent{Text: "msg3"})

	assert.Equal(t, 2, am.GetMessageCount())
	assert.Equal(t, "msg2", am.History[0].Content.Content.(DummyContent).Text)
}

func TestCopy(t *testing.T) {
	am := NewAgentMemory()
	am.InitializeTurn()
	am.AddMessage("user", DummyContent{Text: "copy this"})

	copy := am.Copy()
	assert.Equal(t, am.History, copy.History)
	assert.Equal(t, am.MaxMessages, copy.MaxMessages)
	assert.Equal(t, am.CurrentTurnId, copy.CurrentTurnId)
}

func TestToJsonAndFromJson(t *testing.T) {
	am := NewAgentMemory()
	am.InitializeTurn()
	am.AddMessage("user", DummyContent{Text: "serialize me"})

	jsonStr, err := am.ToJson()
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonStr)

	newAm := NewAgentMemory()
	err = newAm.FromJson(jsonStr)
	assert.NoError(t, err)
	assert.Equal(t, 1, newAm.GetMessageCount())

	// Compare content using marshal to avoid type assertion
	orig, _ := json.Marshal(am)
	deserialized, _ := json.Marshal(newAm)
	assert.JSONEq(t, string(orig), string(deserialized))
}

func TestDeleteMessagesByTurnId_Success(t *testing.T) {
	am := NewAgentMemory()
	am.InitializeTurn()
	turnId := am.CurrentTurnId
	am.AddMessage("user", DummyContent{Text: "delete me"})

	err := am.DeleteMessagesByTurnId(turnId)
	assert.NoError(t, err)
	assert.Equal(t, 0, am.GetMessageCount())
}

func TestDeleteMessagesByTurnId_NotFound(t *testing.T) {
	am := NewAgentMemory()
	am.AddMessage("user", DummyContent{Text: "won't be deleted"})

	err := am.DeleteMessagesByTurnId("non-existent-id")
	assert.Error(t, err)
}
