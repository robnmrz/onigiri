package utils

import "reflect"

// Function that accepts an interface and gets the type name
func GetTypeName(i interface{}) string {
	// Get the type using reflection
	t := reflect.TypeOf(i)

	// Check if it's a pointer and dereference if needed
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Return the type name
	return t.Name()
}
