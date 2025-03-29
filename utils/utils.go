package utils

import "reflect"

// GetTypeName returns the name of a type passed in, safely handling nil
func GetTypeName(i any) string {
	// Handle nil input
	if i == nil {
		return ""
	}

	t := reflect.TypeOf(i)

	// Check if it's a pointer and dereference if needed
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}
