package utils

import (
	"fmt"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// ExtractUser extracts a target user from a reply or command argument
func ExtractUser(b *tb.Bot, m *tb.Message) (*tb.User, error) {
	if m.ReplyTo != nil {
		return m.ReplyTo.Sender, nil
	}

	args := strings.Split(m.Text, " ")
	if len(args) > 1 {
		arg := args[1]

		if id, err := strconv.Atoi(arg); err == nil {
			chat, err := b.ChatByID(fmt.Sprintf("%d", id))
			if err == nil {
				return &tb.User{ID: int64(chat.ID), FirstName: chat.FirstName, Username: chat.Username}, nil
			}
		}

		chat, err := b.ChatByID(arg)
		if err == nil {
			return &tb.User{ID: int64(chat.ID), FirstName: chat.FirstName, Username: chat.Username}, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}
