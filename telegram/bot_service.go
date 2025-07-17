package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"nutritionpro-notify/notification"
)

type BotService struct {
	api      *tgbotapi.BotAPI
	notifier *notification.MealNotifier
	users    map[int64]*UserState
}

func NewBotService(token string) (*BotService, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	notifier := notification.NewMealNotifier(api)

	service := &BotService{
		api:      api,
		notifier: notifier,
		users:    make(map[int64]*UserState),
	}

	return service, nil
}

func (s *BotService) Start() error {
	if err := s.notifier.Start(); err != nil {
		return err
	}

	s.api.Debug = true
	log.Printf("Authorized on account %s", s.api.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := s.api.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			s.handleMessage(update.Message)
		}
	}

	return nil
}

func (s *BotService) Stop() {
	s.notifier.Stop()
}
