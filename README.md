# i18n Library

![GitHub Tag](https://img.shields.io/github/v/tag/go-universal/i18n?sort=semver&label=version)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-universal/i18n.svg)](https://pkg.go.dev/github.com/go-universal/i18n)
[![License](https://img.shields.io/badge/license-ISC-blue.svg)](https://github.com/go-universal/i18n/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-universal/i18n)](https://goreportcard.com/report/github.com/go-universal/i18n)
![Contributors](https://img.shields.io/github/contributors/go-universal/i18n)
![Issues](https://img.shields.io/github/issues/go-universal/i18n)

The `i18n` library provides a simple and extensible solution for handling translations, pluralization, and localization in Go applications. It supports multiple locales, JSON-based translation files, and custom pluralization rules.

## Features

- JSON-based translation management.
- Pluralization support for different languages.
- Customizable localization with placeholders.
- Thread-safe operations for concurrent use.

## Installation

```bash
go get github.com/go-universal/i18n
```

Sample usage:

```go
package main

import (
    "fmt"
    "github.com/go-universal/i18n"
    "golang.org/x/text/language"
)

func main() {
    // Define translations
    en := `{
        "welcome": "Hello {name}, Welcome!",
        "notification": {
            "zero": "No new messages",
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

    // Create a translator
    tr := i18n.NewTranslator("en", language.English)
    tr.AddLocale("fa", &language.Persian)
    tr.LoadBytes("en", []byte(en))
    tr.LoadBytes("fa", []byte(fa))

    // Translate a simple message
    fmt.Println(tr.Translate("en", "welcome", map[string]any{"name": "John"}))
    // Output: Hello John, Welcome!

    // Translate a plural message
    fmt.Println(tr.Plural("fa", "notification", 0, nil))
    // Output: پیام جدیدی ندارید

    fmt.Println(tr.Plural("en", "notification", 5, map[string]any{"count": 5}))
    // Output: You have 5 new messages
}
```

## API Documentation

- `NewTranslator(locale string, formatter language.Tag) Translator`: Creates a new translator with the default locale and formatter.
- `AddLocale(key string, formatter *language.Tag)`: Adds a new locale with the given key and optional formatter.
- `LoadBytes(locale string, content ...[]byte)`: Loads translation JSON content from byte slices for the specified locale.
- `LoadFiles(locale string, files ...string) error`: Loads translation content from JSON files for the specified locale.
- `AddMessage(locale, key, message string, options ...PluralOption)`: Adds a new message with the given key and message string, with optional pluralization rules.
- `Translate(locale, key string, values map[string]any) string`: Translates a message identified by the key using the provided values.
- `Plural(locale, key string, count int, values map[string]any) string`: Translates a plural message identified by the key based on the count and provided values.

## License

This library is licensed under the ISC License. See the [LICENSE](LICENSE) file for details.
