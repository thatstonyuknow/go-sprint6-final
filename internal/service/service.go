package service

import (
	"errors"
	"strings"
	"unicode"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

// Convert automatically determines whether the input string is plain text or Morse code,
// and performs the corresponding conversion. It uses the trimmed input for all operations,
// checks for ambiguous input (e.g. non-text containing invalid Morse symbols),
// and returns an error if the input is empty or ambiguous.
func Convert(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New("input is empty")
	}

	// Check if the input contains letters or digits, indicating plain text.
	isText := strings.ContainsFunc(trimmed, func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	})

	// If not plain text, ensure the input consists solely of valid Morse characters.
	if !isText {
		for _, r := range trimmed {
			if r != '.' && r != '-' && r != ' ' {
				return "", errors.New("input is ambiguous: not plain text or valid Morse code")
			}
		}
	}

	if isText {
		// Convert plain text to Morse code using the trimmed input.
		return morse.ToMorse(trimmed), nil
	}
	// Convert Morse code to plain text using the trimmed input.
	return morse.ToText(trimmed), nil
}
