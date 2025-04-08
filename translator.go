package i18n

import (
	"os"
	"strings"
	"sync"

	"golang.org/x/text/language"
)

// Translator defines an interface for handling translations and localization.
type Translator interface {
	// AddLocale adds a new locale with the given key and formatter.
	AddLocale(key string, formatter *language.Tag)

	// LoadBytes loads translation JSON content from byte slices for the specified locale.
	LoadBytes(locale string, content ...[]byte)

	// LoadFiles loads translation content from JSON files for the specified locale.
	LoadFiles(locale string, files ...string) error

	// AddMessage adds a new message with the given key and message string.
	AddMessage(locale, key, message string, options ...PluralOption)

	// Translate translates a message identified by the key using the provided values.
	Translate(locale, key string, values map[string]any) string

	// Plural translates a plural message identified by the key based on the count and provided values.
	Plural(locale, key string, count int, values map[string]any) string
}

type translator struct {
	defaultLocale string
	defaultTag    language.Tag
	locales       map[string]*localization
	mutex         sync.RWMutex
}

// NewTranslator creates a new translator with the default locale and formatter.
func NewTranslator(locale string, formatter language.Tag) Translator {
	return &translator{
		defaultLocale: locale,
		defaultTag:    formatter,
		locales: map[string]*localization{
			locale: newLocalization(formatter),
		},
	}
}

func (t *translator) AddLocale(key string, formatter *language.Tag) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if formatter == nil {
		t.locales[key] = newLocalization(t.defaultTag)
	} else {
		t.locales[key] = newLocalization(*formatter)
	}
}

func (t *translator) LoadBytes(locale string, contents ...[]byte) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	localeKey, exists := t.resolveLocale(locale)
	if !exists {
		return
	}

	for _, content := range contents {
		t.locales[localeKey].load(content)
	}
}

func (t *translator) LoadFiles(locale string, files ...string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	localeKey, exists := t.resolveLocale(locale)
	if !exists {
		return nil
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		t.locales[localeKey].load(content)
	}
	return nil
}

func (t *translator) AddMessage(locale, key, message string, options ...PluralOption) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	message = strings.TrimSpace(message)
	localeKey, exists := t.resolveLocale(locale)
	if !exists || message == "" {
		return
	}

	t.locales[localeKey].addMessage(key, newMessage(message, options...))
}

func (t *translator) Translate(locale, key string, values map[string]any) string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	localeKey, exists := t.resolveLocale(locale)
	if !exists {
		localeKey, exists = t.resolveLocale("")
		if !exists {
			return ""
		}
	}

	result := t.locales[localeKey].translate(key, -1, values)
	if result == "" && localeKey != t.defaultLocale {
		result = t.locales[t.defaultLocale].translate(key, -1, values)
	}

	return result
}

func (t *translator) Plural(locale, key string, count int, values map[string]any) string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	localeKey, exists := t.resolveLocale(locale)
	if !exists {
		localeKey, exists = t.resolveLocale("")
		if !exists {
			return ""
		}
	}

	result := t.locales[localeKey].translate(key, count, values)
	if result == "" && localeKey != t.defaultLocale {
		result = t.locales[t.defaultLocale].translate(key, count, values)
	}

	return result
}

// resolveLocale resolves the locale key. If the locale is empty, it defaults to the default locale.
func (t *translator) resolveLocale(locale string) (string, bool) {
	if locale == "" {
		locale = t.defaultLocale
	}

	localization, exists := t.locales[locale]
	return locale, exists && localization != nil
}
