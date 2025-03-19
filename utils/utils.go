package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"math/rand"

	"github.com/a-h/templ"

	twmerge "github.com/Oudwins/tailwind-merge-go"
)

// TwMerge combines Tailwind classes and handles conflicts.
// Later classes override earlier ones with the same base.
// Example: "bg-red-500 hover:bg-blue-500", "bg-green-500" → "hover:bg-blue-500 bg-green-500"
func TwMerge(classes ...string) string {
	return twmerge.Merge(classes...)
}

// TwIf returns a class if a condition is true, otherwise an empty string
// Example: "bg-red-500", true → "bg-red-500", false → ""
func TwIf(class string, condition bool) string {
	result := templ.KV(class, condition)
	if result.Value == true {
		return result.Key
	}
	return ""
}

// TwIfElse returns trueClass if condition is true, otherwise falseClass
// Example: true, "bg-red-500", "bg-gray-300" → "bg-red-500", false, "bg-red-500", "bg-gray-300" → "bg-gray-300"
func TwIfElse(condition bool, trueClass string, falseClass string) string {
	if condition {
		return trueClass
	}
	return falseClass
}

// mergeAttributes merges multiple Attributes into one
func MergeAttributes(attrs ...templ.Attributes) templ.Attributes {
	merged := templ.Attributes{}
	for _, attr := range attrs {
		for k, v := range attr {
			merged[k] = v
		}
	}
	return merged
}

// RandomID returns a random ID string
// Example: "id-123456"
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
