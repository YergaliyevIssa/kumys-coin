package main

import (
	"context"
	"fmt"
	"html"
	"kumys-coin/tgbot/pkg/ai"
	"kumys-coin/tgbot/pkg/consts"
	"kumys-coin/tgbot/pkg/doctors"
	"kumys-coin/tgbot/pkg/session"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/dgraph-io/badger/v3"
	tele "gopkg.in/telebot.v4"
)

const (
	BotName     = "SmartDoc"
	LocalDBPath = "db/tgbot"
)

const (
	SectionMainWelcome     = `На что жалуетесь?`
	SectionAnalysisWelcome = `В этой секции вы можете отправить свои анализы (фото, скрины)`
	DoctorsPreText         = `На основе Ваших данных, мы рекомендуем обратиться к терапевту. Вот специалисты из Вашего города, к которым вы можете записаться.`
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
	btnProfile := menu.Text("👤 Профиль")
	btnAnalysis := menu.Text("Анализы")
	btnMain := menu.Text("Главная")

	menu.Reply(
		menu.Row(btnMain),
		menu.Row(btnProfile),
		menu.Row(btnAnalysis),
	)

	// Create a profile menu
	profileMenu := &tele.ReplyMarkup{ResizeKeyboard: true}
	btnProfileChange := menu.Text("Изменить данные")

	profileMenu.Reply(
		menu.Row(btnProfileChange),
		menu.Row(btnMain),
	)

	aiClient := ai.NewClient(os.Getenv("AI_BASE_URL"))

	// Open the BadgerDB database located at dbPath
	opts := badger.DefaultOptions(LocalDBPath).WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sessionRepo := session.NewSessionRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/start", func(c tele.Context) error {
		if err := sessionRepo.CreateSession(&session.Session{
			UserID:    fmt.Sprintf("%d", c.Sender().ID),
			State:     consts.StateInSectionMain,
			ExpiresAt: time.Now().Add(consts.UserSessionTTL),
		}); err != nil {
			return err
		}

		slog.Info("new user session", "userID", c.Sender().ID, "state", consts.StateInSectionMain)

		return c.Send(getWelcomeMessage(), menu)
	})

	// Handle main button
	b.Handle(&btnMain, func(c tele.Context) error {
		if err := c.Send(SectionMainWelcome, menu); err != nil {
			return err
		}

		slog.Info("change user state", "userID", c.Sender().ID, "state", consts.StateInSectionMain)

		return sessionRepo.ChangeUserState(
			fmt.Sprintf("%d", c.Sender().ID),
			consts.StateInSectionMain,
		)
	})

	b.Handle("/testmarkdown", func(c tele.Context) error {
		markdownText := "This is *bold*, _italic_, and `code`.\n" +
			"[Link](https://example.com)\n" +
			"```\nCode block\n```"

		safeText := escapeMarkdown(markdownText)

		return c.Send(safeText, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	})

	// Handle analysis button
	b.Handle(&btnAnalysis, func(c tele.Context) error {
		if err := c.Send(SectionAnalysisWelcome, menu); err != nil {
			return err
		}

		slog.Info("change user state", "userID", c.Sender().ID, "state", consts.StateInSectionAnalysis)

		return sessionRepo.ChangeUserState(
			fmt.Sprintf("%d", c.Sender().ID),
			consts.StateInSectionAnalysis,
		)
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

		return c.Send(profile, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, profileMenu)
	})

	b.Handle("/echo", func(c tele.Context) error {
		text := c.Text()

		slog.Info("got text", "text", text)
		return c.Send(fmt.Sprintf("* привет *"), &tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		session, err := sessionRepo.GetSession(fmt.Sprintf("%d", c.Sender().ID))
		if err != nil {
			return err
		}

		text := c.Text()

		slog.Info("got text", "userID", c.Sender().ID, "text", text)

		switch session.State {
		case consts.StateInSectionMain:
			slog.Info("got text in section main", "userID", c.Sender().ID, "state", session.State)

			resp, err := aiClient.GetDiagnosises(getDefaultContext(), text)
			if err != nil {
				return err
			}

			slog.Info("send diagnoses", "userID", c.Sender().ID, "diagnoses", resp.Diagnosises)
			for _, item := range resp.Diagnosises {
				if err = c.Send(escapeMarkdown(item), menu, &tele.SendOptions{ParseMode: tele.ModeMarkdown}); err != nil {
					slog.Error("send failed", "err", err)
				}
			}

			if err := c.Send(DoctorsPreText); err != nil {
				slog.Error("send doctor pre text", "err", err)
			}

			for _, doctor := range doctors.Doctors {
				slog.Info("photo url", "path", doctor.PhotoURL)
				// Create a photo from a URL
				photo := &tele.Photo{File: tele.FromDisk(doctor.PhotoURL)}

				// Send the photo with a caption
				if err := c.Send(&tele.Photo{
					File: photo.File,
				}); err != nil {
					slog.Error("send photo", "err", err)
				}

				if err := c.Send(doctor.String()); err != nil {
					slog.Error("send doctor info", "err", err)
				}
			}

			return nil
		case consts.StateChangingProfile:
			//
			return nil
		}

		return c.Send("unexpected state", menu)
	})

	b.Handle(tele.OnPhoto, func(c tele.Context) error {
		session, err := sessionRepo.GetSession(fmt.Sprintf("%d", c.Sender().ID))
		if err != nil {
			return err
		}

		photo := c.Message().Photo

		slog.Info("got photo", "userID", c.Sender().ID, "photo size", photo.FileSize)

		switch session.State {
		case consts.StateInSectionAnalysis:
			slog.Info("got text in section analysis", "userID", c.Sender().ID, "state", session.State)

			file, err := b.File(&photo.File)
			if err != nil {
				return err
			}

			resp, err := aiClient.SendAnalysis(getDefaultContext(), file)
			if err != nil {
				return err
			}

			slog.Info("send analysis", "userID", c.Sender().ID, "analytics", resp.Analytics)
			return c.Send(escapeMarkdown(resp.Analytics), menu, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
		}

		return c.Send("...", menu)
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

func escapeMarkdown(text string) string {
	return html.EscapeString(text)
}
