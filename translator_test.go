package i18n_test

import (
	"testing"

	"github.com/go-universal/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestTranslator(t *testing.T) {
	// Raw data
	en := `{
		"welcome": "Hello {name}, Welcome!",
		"notification": {
			"zero": "No new message",
			"one": "You have a new message",
			"other": "You have {count} new messages"
		}
	}`
	fa := `{
		"welcome": "سلام {name}, خوش آمدید!",
		"notification": {
			"zero": "پیام جدیدی ندارید",
			"one": "یک پیام جدید دارید",
			"other": "{count} پیام جدید دارید"
		}
	}`

	// Create translator
	tr := i18n.NewTranslator("en", language.English)
	tr.AddLocale("fa", &language.Persian)
	tr.LoadBytes("en", []byte(en))
	tr.LoadBytes("fa", []byte(fa))

	// Run tests
	t.Run("Simple", func(t *testing.T) {
		res := tr.Translate("", "welcome", map[string]any{"name": "John doe"})
		assert.Equal(t, "Hello John doe, Welcome!", res, "Translation mismatch for 'welcome'")
	})

	t.Run("Plural", func(t *testing.T) {
		res := tr.Plural("fa", "notification", 0, map[string]any{"name": "John doe"})
		assert.Equal(t, "پیام جدیدی ندارید", res, "Translation mismatch for 'notification' with count 0 in 'fa'")

		res = tr.Plural("ar", "notification", 10, map[string]any{"count": 10})
		assert.Equal(t, "You have 10 new messages", res, "Translation mismatch for 'notification' with count 10 in 'ar'")
	})
}
