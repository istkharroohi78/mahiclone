package misc

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterWatcherHandlers(b *tb.Bot) {
	// Telebot v2 doesn't have native VoiceChat events.
	// This usually requires parsing raw updates or using a specialized MTProto lib.

	// Mock implementation to represent the logic:
	// If Voice Chat ends -> Stop Stream Forcefully
	/*
		b.Handle(tb.OnVoiceChatEnded, func(m *tb.Message) {
			stream.ClearQueue(m.Chat.ID)
			// Stop core PyTgCalls
		})
	*/
}
