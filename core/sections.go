package core

import (
	"fmt"
	"strings"
)

// Helper functions for Telegram Markdown formatting
func Bold(x string) string {
	return fmt.Sprintf("**%s:** ", x)
}

func BoldUl(x string) string {
	return fmt.Sprintf("**--%s:**-- ", x)
}

func Mono(x interface{}) string {
	return fmt.Sprintf("`%v`\n", x)
}

// Section generates a formatted string for a section with key-value pairs.
func Section(title string, body map[string]interface{}, indent int, underline bool) string {
	var text strings.Builder
	n := "\n"
	w := " "

	if underline {
		text.WriteString(BoldUl(title) + n)
	} else {
		text.WriteString(Bold(title) + n)
	}

	indentStr := strings.Repeat(w, indent)

	for key, value := range body {
		if value != nil {
			text.WriteString(indentStr)
			text.WriteString(Bold(key))

			// Type assertion to check if value is a list of strings
			if strSlice, ok := value.([]string); ok && len(strSlice) > 0 {
				text.WriteString(strSlice[0] + n)
			} else {
				text.WriteString(Mono(value))
			}
		}
	}

	return text.String()
}
