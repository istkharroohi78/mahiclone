package inline

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func StatsButtons(loc map[string]string, status bool) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}

	var topRow []tb.Btn
	if status {
		topRow = []tb.Btn{
			CreateBtn(m, loc["SA_B_2"], "bot_stats_sudo", "", Primary, 0, false),
			CreateBtn(m, loc["SA_B_3"], "TopOverall", "", Primary, 0, false),
		}
	} else {
		topRow = []tb.Btn{
			CreateBtn(m, loc["SA_B_1"], "TopOverall", "", Primary, 0, false),
		}
	}

	m.Inline(
		m.Row(topRow...),
		m.Row(CreateBtn(m, loc["CLOSE_BUTTON"], "close", "", Danger, 0, false)),
	)
	return m
}

func BackStatsButtons(loc map[string]string) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	m.Inline(
		m.Row(
			CreateBtn(m, loc["BACK_BUTTON"], "stats_back", "", Primary, 0, false),
			CreateBtn(m, loc["CLOSE_BUTTON"], "close", "", Danger, 0, false),
		),
	)
	return m
}
