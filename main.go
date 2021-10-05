package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zilstream/bot/helpers"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

type server struct {
	db     *gorm.DB
	router *echo.Echo
	bot    *tb.Bot
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

	s := server{
		db:     SetupDatabase(),
		router: echo.New(),
		bot:    b,
	}

	s.RegisterBotHandlers()
	s.bot.Start()

	s.router.Logger.Fatal(s.router.Start(":8080"))
}
