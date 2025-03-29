package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Dummy struct for testing
type MyStruct struct {
	Name string
}

func TestGetTypeName_BasicType(t *testing.T) {
	name := GetTypeName(123)
	assert.Equal(t, "int", name)

	name = GetTypeName("hello")
	assert.Equal(t, "string", name)
}

func TestGetTypeName_Struct(t *testing.T) {
	s := MyStruct{Name: "test"}
	name := GetTypeName(s)
	assert.Equal(t, "MyStruct", name)
}

func TestGetTypeName_PointerToStruct(t *testing.T) {
	s := &MyStruct{Name: "test"}
	name := GetTypeName(s)
	assert.Equal(t, "MyStruct", name)
}

func TestGetTypeName_AnonymousStruct(t *testing.T) {
	anon := struct {
		ID int
	}{ID: 1}

	name := GetTypeName(anon)
	assert.Equal(t, "", name) // anonymous structs return an empty string
}

func TestGetTypeName_Nil(t *testing.T) {
	var x any
	name := GetTypeName(x)
	assert.Equal(t, "", name)
}
