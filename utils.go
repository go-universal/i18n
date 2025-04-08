package i18n

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tidwall/gjson"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// toString converts a value `v` to its string representation based on the specified language `l`.
// It handles various types, including pointers, integers, floats, and types implementing fmt.Stringer.
func toString(l language.Tag, v any) string {
	// Handle pointer types by dereferencing
	if val := reflect.ValueOf(v); val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ""
		}

		v = val.Elem().Interface()
	}

	// Convert value to string based on its type
	switch v := v.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return message.NewPrinter(l).Sprintf("%d", v)
	case float32, float64:
		return message.NewPrinter(l).Sprintf("%.2f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// translateJson translates a JSON result based on the given locale, count, and placeholder values.
// It selects the appropriate translation string based on the count and replaces placeholders with their values.
func translateJson(locale language.Tag, json gjson.Result, count int, values map[string]any) string {
	// Determine the appropriate message based on the count
	var msg string
	switch {
	case count == 0:
		msg = json.Get("zero").Str
	case count == 1:
		msg = json.Get("one").Str
	case count == 2:
		msg = json.Get("two").Str
	case count > 2 && count <= 10:
		msg = json.Get("few").Str
	default:
		msg = json.Get("many").Str
	}

	// Fallback to "other" or the entire JSON string if no specific message is found
	if msg == "" {
		msg = json.Get("other").Str
		if msg == "" {
			msg = json.Str
		}
	}

	// Replace placeholders in the message with corresponding values
	for k, v := range values {
		placeholder := "{" + k + "}"
		msg = strings.ReplaceAll(msg, placeholder, toString(locale, v))
	}

	return msg
}
