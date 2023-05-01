package tg_engine

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"time"
	"webcamer/webcamer"
)

var (
	invalidTokenError    = errors.New("invalid token string")
	invalidWebcamerError = errors.New("invalid webcamer")
)

type Engine struct {
	adminsMap map[int64]struct{}
	api       *tgbotapi.BotAPI
	webcamer  webcamer.WebCamerer
}

func NewEngine(token string, wcamer webcamer.WebCamerer, admins ...int64) (*Engine, error) {
	if len(token) == 0 {
		return nil, invalidTokenError
	}
	if wcamer == nil {
		return nil, invalidWebcamerError
	}
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	e := &Engine{
		adminsMap: make(map[int64]struct{}, len(admins)),
		webcamer:  wcamer,
		api:       api,
	}

	for _, id := range admins {
		if _, ok := e.adminsMap[id]; !ok {
			e.adminsMap[id] = struct{}{}
		}
	}
	return e, nil
}

func (e *Engine) Run() {
	//TODO: refactor this method
	u := tgbotapi.NewUpdate(0)

	updates := e.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			logrus.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.From != nil {
				if _, ok := e.adminsMap[update.Message.From.ID]; ok {
					switch update.Message.Text {
					case "/snapshot":
						photo, err := e.webcamer.DoSnapshot()
						if err != nil {
							_, _ = e.api.Send(
								tgbotapi.NewMessage(
									update.Message.Chat.ID,
									fmt.Sprintf("Не удалось сделать фото: %v", err),
								),
							)
							continue
						}
						msg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(photo))
						msg.ReplyToMessageID = update.Message.MessageID
						msg.Caption = fmt.Sprintf("%v", time.Now().Format("2006-01-02 15:04:05"))

						_, err = e.api.Send(msg)
						if err != nil {
							_, _ = e.api.Send(
								tgbotapi.NewMessage(
									update.Message.Chat.ID,
									fmt.Sprintf("Не удалось сделать фото: %v", err),
								),
							)
							logrus.Error(err)
						}
					case "/video":
						_, _ = e.api.Send(
							tgbotapi.NewMessage(
								update.Message.Chat.ID,
								"Начинаю запись видео, дождитесь ответа 🙂",
							),
						)
						video, err := e.webcamer.DoVideo(30)
						if err != nil {
							_, _ = e.api.Send(
								tgbotapi.NewMessage(
									update.Message.Chat.ID,
									fmt.Sprintf("Не удалось сделать фото: %v", err),
								),
							)
							continue
						}
						msg := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath(video))
						msg.ReplyToMessageID = update.Message.MessageID
						msg.Caption = fmt.Sprintf("%v", time.Now().Format("2006-01-02 15:04:05"))

						_, err = e.api.Send(msg)
						if err != nil {
							_, _ = e.api.Send(
								tgbotapi.NewMessage(
									update.Message.Chat.ID,
									fmt.Sprintf("Не удалось сделать фото: %v", err),
								),
							)
							logrus.Error(err)
						}
					}

				}
			}
		}
	}
}

func (e *Engine) Stop() {
	e.api.StopReceivingUpdates()
}
