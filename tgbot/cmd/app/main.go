package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	tele "gopkg.in/telebot.v4"
)

const (
	BotName = "SuperAppteka"
)

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send(getWelcomeMessage())
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		// All the text messages that weren't
		// captured by existing handlers.
		text := c.Text()

		return c.Send(text)
	})

	slog.Info("starting tgbot")
	b.Start()
}

func getWelcomeMessage() string {
	return fmt.Sprintf("Добро пожаловать в %s!\n"+
		`Здесь вы можете найти рекомендации по лечению различных заболеваний, а также полезную информацию о таблетках и лекарствах. Просто задайте свой вопрос, и я помогу вам разобраться!`,
		BotName)
}
