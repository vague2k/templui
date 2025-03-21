package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"math/rand"

	"github.com/a-h/templ"

	twmerge "github.com/Oudwins/tailwind-merge-go"
)

// TwMerge combines Tailwind classes and resolves conflicts.
// Example: "bg-red-500 hover:bg-blue-500", "bg-green-500" → "hover:bg-blue-500 bg-green-500"
func TwMerge(classes ...string) string {
	return twmerge.Merge(classes...)
}

// TwIf returns value if condition is true, otherwise an empty value of type T.
// Example: true, "bg-red-500" → "bg-red-500"
func If[T comparable](condition bool, value T) T {
	var empty T
	if condition {
		return value
	}
	return empty
}

// TwIfElse returns trueValue if condition is true, otherwise falseValue.
// Example: true, "bg-red-500", "bg-gray-300" → "bg-red-500"
func IfElse[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// MergeAttributes combines multiple Attributes into one.
// Example: MergeAttributes(attr1, attr2) → combined attributes
func MergeAttributes(attrs ...templ.Attributes) templ.Attributes {
	merged := templ.Attributes{}
	for _, attr := range attrs {
		for k, v := range attr {
			merged[k] = v
		}
	}
	return merged
}

// RandomID generates a random ID string.
// Example: RandomID() → "id-123456"
func RandomID() string {
	return fmt.Sprintf("id-%d", rand.Intn(1000000))
}

// create a ifelse func withb bool prop but return any
func IfElse(b bool, t, f any) any {
	if b {
		return t
	}
	return f
}

// float to string create one with standard lib
func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// int to string create one with standard lib
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// json stringify create one with standard lib
func JSONStringify(v any) string {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "[]" // Rückgabe eines leeren Arrays/Objekts bei Fehler
	}
	return string(jsonBytes)
}

// Ptr returns a pointer to the given value of any type.
func Ptr[T any](v T) *T {
	return &v
}
