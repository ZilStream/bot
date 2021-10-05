package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/zilstream/bot/models"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *server) RegisterBotHandlers() {
	s.bot.Handle("/start", func(m *tb.Message) {
		text := "Hi there!\n\nTo retrieve the current price for a token:\n/zs [symbol]\n/zs gzil\n/zs zwap\n/zs port\netc.\n\nTo retrieve the price with more detail use /zss.\n\nType /help to see this message again."
		s.bot.Send(m.Sender, text, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	})

	s.bot.Handle("/help", func(m *tb.Message) {
		text := "Hi there!\n\nTo retrieve the current price for a token:\n/zs [symbol]\n/zs gzil\n/zs zwap\n/zs port\netc.\n\nTo retrieve the price with more detail use /zss."
		s.bot.Send(m.Sender, text, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	})

	s.bot.Handle("/zs", func(m *tb.Message) {
		symbol := m.Payload

		res, err := http.Get(fmt.Sprintf("https://api.zilstream.com/tokens/%s", strings.ToLower(symbol)))
		if err != nil {
			s.bot.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			s.bot.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		var token models.TokenDetail
		err = json.Unmarshal(body, &token)
		if err != nil {
			s.bot.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		if token.Symbol == "" {
			s.bot.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		p := message.NewPrinter(language.English)

		text := p.Sprintf("<b>%s (%s)</b>\n<b>%.2f ZIL - $%.2f</b>\n<pre>Change (24h):  %.2f%%\nMarket Cap:    $%.2f\nVolume (24h):  $%.2f\nCirc. Supply:  %.0f</pre>\n<a href='https://zilstream.com/tokens/%s'>View %s on ZilStream</a>",
			token.Name,
			token.Symbol,
			token.Rate,
			token.RateUSD,
			token.MarketData.ChangePercentage24H,
			token.MarketCap,
			token.DailyVolume,
			token.CurrentSupply,
			strings.ToLower(token.Symbol),
			token.Symbol,
		)

		if !token.Listed {
			text = p.Sprintf("<b>%s (%s)</b>\n<b>UNLISTED: BE EXTRA CAUTIOUS</b>\n%.2f ZIL - $%.2f\n<pre>Change (24h):  %.2f%%\nMarket Cap:    $%.2f\nVolume (24h):  $%.2f\nCirc. Supply:  %.0f</pre>\n<a href='https://zilstream.com/tokens/%s'>View %s on ZilStream</a>",
				token.Name,
				token.Symbol,
				token.Rate,
				token.RateUSD,
				token.MarketData.ChangePercentage24H,
				token.MarketCap,
				token.DailyVolume,
				token.CurrentSupply,
				strings.ToLower(token.Symbol),
				token.Symbol,
			)
		}

		_, err = s.bot.Send(m.Chat, text, &tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeHTML})
		if err != nil {
			fmt.Println(err)
		}
	})

	s.bot.Handle("/zss", func(m *tb.Message) {
		symbol := m.Payload

		res, err := http.Get(fmt.Sprintf("https://api.zilstream.com/tokens/%s", strings.ToLower(symbol)))
		if err != nil {
			s.bot.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			s.bot.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		var token models.TokenDetail
		err = json.Unmarshal(body, &token)
		if err != nil {
			s.bot.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		if token.Symbol == "" {
			s.bot.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		p := message.NewPrinter(language.English)
		text := p.Sprintf("<b>%s (%s)</b>\n<b>%.2f ZIL - $%.2f</b>\n<pre>ATH:           %.2f\nChange (24h):  %.2f%%\nChange (7d):   %.2f%%\nMarket Cap:    $%.2f\nFully Diluted: $%.2f\nVolume (24h):  $%.2f\nCirc. Supply:  %.0f\nLiquidity:     $%.2f \n               %.2f ZIL\n               %.2f %s</pre>\n<a href='https://zilstream.com/tokens/%s'>View %s on ZilStream</a>",
			token.Name,
			token.Symbol,
			token.Rate,
			token.RateUSD,
			token.MarketData.ATH,
			token.MarketData.ChangePercentage24H,
			token.MarketData.ChangePercentage7D,
			token.MarketCap,
			token.MarketData.FullyDilutedValuation,
			token.DailyVolume,
			token.CurrentSupply,
			token.MarketData.CurrentLiquidity,
			token.MarketData.ZilReserve,
			token.MarketData.TokenReserve,
			token.Symbol,
			strings.ToLower(token.Symbol),
			token.Symbol,
		)

		if !token.Listed {
			text = p.Sprintf("<b>%s (%s)</b>\n<b>UNLISTED: BE EXTRA CAUTIOUS</b>\n%.2f ZIL - $%.2f\n<pre>ATH:           %.2f\nChange (24h):  %.2f%%\nChange (7d):   %.2f%%\nMarket Cap:    $%.2f\nFully Diluted: $%.2f\nVolume (24h):  $%.2f\nCirc. Supply:  %.0f\nLiquidity:     $%.2f \n               %.2f ZIL\n               %.2f %s</pre>\n<a href='https://zilstream.com/tokens/%s'>View %s on ZilStream</a>",
				token.Name,
				token.Symbol,
				token.Rate,
				token.RateUSD,
				token.MarketData.ATH,
				token.MarketData.ChangePercentage24H,
				token.MarketData.ChangePercentage7D,
				token.MarketCap,
				token.MarketData.FullyDilutedValuation,
				token.DailyVolume,
				token.CurrentSupply,
				token.MarketData.CurrentLiquidity,
				token.MarketData.ZilReserve,
				token.MarketData.TokenReserve,
				token.Symbol,
				strings.ToLower(token.Symbol),
				token.Symbol,
			)
		}

		_, err = s.bot.Send(m.Chat, text, &tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeHTML})
		if err != nil {
			fmt.Println(err)
		}
	})
}
