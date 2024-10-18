package main

import (
	"context"
	"fmt"
	"kumys-coin/tgbot/pkg/ai"
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

	aiClient := ai.NewClient(os.Getenv("AI_BASE_URL"))

	b.Handle("/start", func(c tele.Context) error {
		return c.Send(getWelcomeMessage())
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		text := c.Text()

		resp, err := aiClient.GetDiagnosises(getDefaultContext(), text)
		if err != nil {
			return err
		}

		for _, item := range resp.Recommendations {
			if err = c.Send(item); err != nil {
				slog.Error("send failed", "err", err)
			}
		}

		return nil
	})

	b.Handle(tele.OnPhoto, func(c tele.Context) error {
		photo := c.Message().Photo

		file, err := b.File(&photo.File)
		if err != nil {
			return err
		}

		resp, err := aiClient.SendAnalysis(getDefaultContext(), file)
		if err != nil {
			return err
		}

		return c.Send(resp.Text)
	})

	slog.Info("starting tgbot")
	b.Start()
}

func getWelcomeMessage() string {
	return fmt.Sprintf("Добро пожаловать в %s!\n"+
		`Здесь вы можете найти рекомендации по лечению различных заболеваний, а также полезную информацию о таблетках и лекарствах. Просто задайте свой вопрос, и я помогу вам разобраться!`,
		BotName)
}

func getDefaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	return ctx
}
