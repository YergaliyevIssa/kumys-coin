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

	// Create a main menu with buttons
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	btnProfile := menu.Text("👤 Profile")
	btnAnalysis := menu.Text("Check analysis")

	menu.Reply(
		menu.Row(btnProfile),
		menu.Row(btnAnalysis),
	)

	aiClient := ai.NewClient(os.Getenv("AI_BASE_URL"))

	b.Handle("/start", func(c tele.Context) error {
		return c.Send(getWelcomeMessage(), menu)
	})

	// Handle Profile button
	b.Handle(&btnProfile, func(c tele.Context) error {
		user := c.Sender()
		profile := fmt.Sprintf("👤 *Profile Information*\n\n"+
			"Name: %s\n"+
			"Username: @%s\n"+
			"User ID: %d\n"+
			"Language Code: %s",
			user.FirstName+" "+user.LastName,
			user.Username,
			user.ID,
			user.LanguageCode)

		return c.Send(profile, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, menu)
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		text := c.Text()

		resp, err := aiClient.GetDiagnosises(getDefaultContext(), text)
		if err != nil {
			return err
		}

		for _, item := range resp.Diagnosises {
			if err = c.Send(item); err != nil {
				slog.Error("send failed", "err", err)
			}
		}

		return c.Send("...", menu)
	})

	// Handle analysis
	b.Handle(&btnAnalysis, func(c tele.Context) error {
		photo := c.Message().Photo

		file, err := b.File(&photo.File)
		if err != nil {
			return err
		}

		resp, err := aiClient.SendAnalysis(getDefaultContext(), file)
		if err != nil {
			return fmt.Errorf("send analysis: %w", err)
		}

		return c.Send(resp.Analytics)
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
