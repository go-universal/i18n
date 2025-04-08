package i18n

import (
	"github.com/tidwall/gjson"
	"golang.org/x/text/language"
)

// localization represents a localization context with a specific language tag,
// a collection of parsed JSON files, and a map of plural messages.
type localization struct {
	tag      language.Tag
	files    []gjson.Result
	messages map[string]pluralMessage
}

// newLocalization creates and initializes a new localization context
// with the specified language tag.
func newLocalization(tag language.Tag) *localization {
	return &localization{
		tag:      tag,
		files:    make([]gjson.Result, 0),
		messages: make(map[string]pluralMessage),
	}
}

// load parses the provided JSON content and appends it to the localization files
// if the content is valid JSON. Invalid JSON is ignored.
func (l *localization) load(content []byte) {
	if gjson.ValidBytes(content) {
		l.files = append(l.files, gjson.ParseBytes(content))
	}
}

// addMessage adds a plural message to the localization context.
// This allows for custom pluralization logic for specific keys.
func (l *localization) addMessage(key string, message pluralMessage) {
	l.messages[key] = message
}

// translate retrieves the localized message for the given key, count, and values.
// It first checks the custom messages map, then searches the JSON files.
// If no match is found, it returns an empty string.
func (l *localization) translate(key string, count int, values map[string]any) string {
	// Check if the key exists in the custom messages map.
	if message, exists := l.messages[key]; exists {
		return message.translate(l.tag, count, values)
	}

	// Search for the key in the JSON files.
	for _, file := range l.files {
		if value := file.Get(key); value.Exists() {
			return translateJson(l.tag, value, count, values)
		}
	}

	// Return an empty string if no translation is found.
	return ""
}
