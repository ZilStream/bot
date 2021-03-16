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

type TokenDetail struct {
	Name          string  `json:"name"`
	Symbol        string  `json:"symbol"`
	CurrentSupply float64 `json:"current_supply"`
	DailyVolume   float64 `json:"daily_volume"`
	MarketCap     float64 `json:"market_cap"`
	Rate          float64 `json:"rate"`
	RateUSD       float64 `json:"rate_usd"`
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

	b.Handle("/zs", func(m *tb.Message) {
		symbol := m.Payload

		res, err := http.Get(fmt.Sprintf("https://api.zilstream.com/tokens/%s", strings.ToLower(symbol)))
		if err != nil {
			b.Send(m.Sender, fmt.Sprintf("Couldn't find %s", symbol))
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			b.Send(m.Sender, fmt.Sprintf("Couldn't find %s", symbol))
		}

		var token TokenDetail
		err = json.Unmarshal(body, &token)
		if err != nil {
			b.Send(m.Sender, fmt.Sprintf("Couldn't find %s", symbol))
		}

		p := message.NewPrinter(language.English)
		text := p.Sprintf("*%s (%s)* \n*%.2f ZIL - $%.2f*\nMarket Cap: $%.2f\nVolume (24h): $%.2f\nCirc. Supply: %.0f\n[View %s on ZilStream](https://zilstream.com/tokens/%s)", token.Name, token.Symbol, token.Rate, token.RateUSD, token.MarketCap, token.DailyVolume, token.CurrentSupply, token.Symbol, strings.ToLower(token.Symbol))
		_, err = b.Send(m.Sender, text, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
		if err != nil {
			fmt.Println(err)
		}
	})

	b.Start()
}
