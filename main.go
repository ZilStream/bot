package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/zilstream/bot/helpers"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	tb "gopkg.in/tucnak/telebot.v2"
)

type MarketData struct {
	ATH                   float64 `json:"ath"`
	ATL                   float64 `json:"atl"`
	Change24H             float64 `json:"change_24h"`
	ChangePercentage24H   float64 `json:"change_percentage_24h"`
	ChangePercentage7D    float64 `json:"change_percentage_7d"`
	ChangePercentage14D   float64 `json:"change_percentage_14d"`
	ChangePercentage30D   float64 `json:"change_percentage_30d"`
	InitSupply            float64 `json:"init_supply"`
	MaxSupply             float64 `json:"max_supply"`
	TotalSupply           float64 `json:"total_supply"`
	CurrentSupply         float64 `json:"current_supply"`
	DailyVolume           float64 `json:"daily_volume"`
	MarketCap             float64 `json:"market_cap"`
	FullyDilutedValuation float64 `json:"fully_diluted_valuation"`
	CurrentLiquidity      float64 `json:"current_liquidity"`
	ZilReserve            float64 `json:"zil_reserve"`
	TokenReserve          float64 `json:"token_reserve"`
}

type TokenDetail struct {
	Name                string     `json:"name"`
	Symbol              string     `json:"symbol"`
	AddressBech32       string     `json:"address_bech32"`
	Icon                string     `json:"icon"`
	Decimals            float64    `json:"decimals"`
	Website             string     `json:"website"`
	Whitepaper          string     `json:"whitepaper"`
	CurrentSupply       float64    `json:"current_supply"`
	DailyVolume         float64    `json:"daily_volume"`
	SupplySkipAddresses string     `json:"supply_skip_addresses"`
	MarketCap           float64    `json:"market_cap"`
	Rate                float64    `json:"rate"`
	RateUSD             float64    `json:"rate_usd"`
	MarketData          MarketData `json:"market_data"`
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  helpers.GetEnv("TG_TOKEN", ""),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(m *tb.Message) {
		text := "Hi there!\n\nTo retrieve the current price for a token:\n/zs [symbol]\n/zs gzil\n/zs zwap\n/zs port\netc.\n\nTo retrieve the price with more detail use /zss.\n\nType /help to see this message again."
		b.Send(m.Sender, text, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	})

	b.Handle("/help", func(m *tb.Message) {
		text := "Hi there!\n\nTo retrieve the current price for a token:\n/zs [symbol]\n/zs gzil\n/zs zwap\n/zs port\netc.\n\nTo retrieve the price with more detail use /zss."
		b.Send(m.Sender, text, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	})

	b.Handle("/zs", func(m *tb.Message) {
		symbol := m.Payload

		res, err := http.Get(fmt.Sprintf("https://api.zilstream.com/tokens/%s", strings.ToLower(symbol)))
		if err != nil {
			b.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			b.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		var token TokenDetail
		err = json.Unmarshal(body, &token)
		if err != nil {
			b.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		if token.Symbol == "" {
			b.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		p := message.NewPrinter(language.English)
		text := p.Sprintf("*%s (%s)* \n*%.2f ZIL - $%.2f*\nMarket Cap: $%.2f\nVolume (24h): $%.2f\nCirc. Supply: %.0f\n[View %s on ZilStream](https://zilstream.com/tokens/%s)", token.Name, token.Symbol, token.Rate, token.RateUSD, token.MarketCap, token.DailyVolume, token.CurrentSupply, token.Symbol, strings.ToLower(token.Symbol))
		_, err = b.Send(m.Chat, text, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
		if err != nil {
			fmt.Println(err)
		}
	})

	b.Handle("/zss", func(m *tb.Message) {
		symbol := m.Payload

		res, err := http.Get(fmt.Sprintf("https://api.zilstream.com/tokens/%s", strings.ToLower(symbol)))
		if err != nil {
			b.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			b.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		var token TokenDetail
		err = json.Unmarshal(body, &token)
		if err != nil {
			b.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		if token.Symbol == "" {
			b.Send(m.Sender, fmt.Sprintf("Couldn't find %s, type /help for more information.", symbol))
			return
		}

		p := message.NewPrinter(language.English)
		text := p.Sprintf("<b>%s (%s) - %.2f ZIL - $%.2f</b>\n<pre>ATH:           %.2f\nChange (24h):  %.2f%%\nChange (7d):   %.2f%%\nMarket Cap:    $%.2f\nFully Diluted: $%.2f\nVolume (24h):  $%.2f\nCirc. Supply:  %.0f\nLiquidity:     $%.2f \n               %.2f ZIL\n               %.2f %s</pre>\n<a href='https://zilstream.com/tokens/%s'>View %s on ZilStream</a>",
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
			token.Symbol,
			strings.ToLower(token.Symbol),
		)
		_, err = b.Send(m.Chat, text, &tb.SendOptions{ParseMode: tb.ModeHTML})
		if err != nil {
			fmt.Println(err)
		}
	})

	b.Start()
}
