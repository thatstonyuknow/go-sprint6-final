package service

import (
	"errors"
	"strings"
	"unicode"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

// Convert automatically determines whether the input string is plain text or Morse code.
// If plain text is provided, the function converts it to Morse code and returns the result.
// Conversely, if Morse code is provided, it converts it to plain text and returns it.
// It utilizes the strings package for parsing and the morse package for conversion.
// If any error occurs (e.g., an unknown character is encountered), an error is returned.
func Convert(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New("input is empty")
	}

	// Determine if the string is Morse code.
	// If the string consists only of dots, dashes, and spaces, consider it as Morse code.
	isMorse := true
	for _, r := range trimmed {
		if r != '.' && r != '-' && !unicode.IsSpace(r) {
			isMorse = false
			break
		}
	}

	if isMorse {
		// Convert from Morse code to plain text.
		result := morse.ToText(trimmed)
		if result == "" {
			return "", errors.New("conversion error: Morse to text resulted in an empty string")
		}
		return result, nil
	} else {
		// Before converting from plain text to Morse code,
		// check that each character (except spaces) has a corresponding Morse representation.
		for _, r := range trimmed {
			if unicode.IsSpace(r) {
				continue
			}
			// Convert character to uppercase because the encoding map stores letters in uppercase.
			if _, ok := morse.DefaultMorse[unicode.ToUpper(r)]; !ok {
				return "", errors.New("conversion error: unknown character " + string(r))
			}
		}
		// Convert from plain text to Morse code.
		result := morse.ToMorse(trimmed)
		if result == "" {
			return "", errors.New("conversion error: text to Morse resulted in an empty string")
		}
		return result, nil
	}
}
