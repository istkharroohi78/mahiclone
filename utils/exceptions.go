package utils

import "fmt"

type AssistantErr struct {
	Message string
}

func NewAssistantErr(msg string) error {
	return &AssistantErr{Message: msg}
}

func (e *AssistantErr) Error() string {
	return fmt.Sprintf("Assistant Error: %s", e.Message)
}
