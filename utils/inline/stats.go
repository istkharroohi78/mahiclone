package inline

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func StatsButtons(loc map[string]string, status bool) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}

	var topRow []tb.Btn
	if status {
		topRow = []tb.Btn{
			createBtn(m, loc["SA_B_2"], "bot_stats_sudo", "", Primary, "", false),
			createBtn(m, loc["SA_B_3"], "TopOverall", "", Primary, "", false),
		}
	} else {
		topRow = []tb.Btn{
			createBtn(m, loc["SA_B_1"], "TopOverall", "", Primary, "", false),
		}
	}

	m.Inline(
		m.Row(topRow...),
		m.Row(createBtn(m, loc["CLOSE_BUTTON"], "close", "", Danger, "", false)),
	)
	return m
}

func BackStatsButtons(loc map[string]string) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	m.Inline(
		m.Row(
			createBtn(m, loc["BACK_BUTTON"], "stats_back", "", Primary, "", false),
			createBtn(m, loc["CLOSE_BUTTON"], "close", "", Danger, "", false),
		),
	)
	return m
}
