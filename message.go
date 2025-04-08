package i18n

import (
	"strings"

	"golang.org/x/text/language"
)

// pluralMessage represents a message with different plural forms.
type pluralMessage struct {
	zero  string // Message for count = 0.
	one   string // Message for count = 1.
	two   string // Message for count = 2.
	few   string // Message for a few counts (e.g., 3-10).
	many  string // Message for many counts (e.g., >10).
	other string // Default message for other counts or fallback.
}

// PluralOption defines a function type for modifying a pluralMessage.
type PluralOption func(*pluralMessage)

// PluralZero sets the message for count = 0.
func PluralZero(msg string) PluralOption {
	return func(m *pluralMessage) { m.zero = msg }
}

// PluralOne sets the message for count = 1.
func PluralOne(msg string) PluralOption {
	return func(m *pluralMessage) { m.one = msg }
}

// PluralTwo sets the message for count = 2.
func PluralTwo(msg string) PluralOption {
	return func(m *pluralMessage) { m.two = msg }
}

// PluralFew sets the message for a few counts (e.g., 3-10).
func PluralFew(msg string) PluralOption {
	return func(m *pluralMessage) { m.few = msg }
}

// PluralMany sets the message for many counts (e.g., >10).
func PluralMany(msg string) PluralOption {
	return func(m *pluralMessage) { m.many = msg }
}

// resolve determines the appropriate message based on the count.
func (m pluralMessage) resolve(count int) string {
	switch {
	case count == 0:
		return m.zero
	case count == 1:
		return m.one
	case count == 2:
		return m.two
	case count > 2 && count <= 10:
		return m.few
	case count > 10:
		return m.many
	default:
		return m.other // Fallback to the default message.
	}
}

// translate returns the resolved message with placeholders replaced by values.
func (m pluralMessage) translate(locale language.Tag, count int, values map[string]any) string {
	message := m.resolve(count)
	if message == "" {
		return ""
	}

	// Replace placeholders (e.g., {key}) with corresponding values.
	for k, v := range values {
		message = strings.ReplaceAll(message, "{"+k+"}", toString(locale, v))
	}
	return message
}

// newMessage creates a new pluralMessage with a default message and optional plural forms.
func newMessage(defaultMessage string, options ...PluralOption) pluralMessage {
	m := &pluralMessage{other: defaultMessage}
	for _, opt := range options {
		opt(m)
	}
	return *m
}
