package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"nutritionpro-notify/nutritionpro"
	"regexp"
	"strings"
)

type UserState struct {
	Stage    int
	ChatID   int64
	PhoneNum string
}

func isValidPhoneNumber(phone string) bool {
	// Czech phone number validation
	// Format: +420 followed by 9 digits
	czechPhoneRegex := regexp.MustCompile(`^\+420[0-9]{9}$`)
	return czechPhoneRegex.MatchString(phone)
}

func (s *BotService) handleMenuForToday(chatID int64) {
	state, exists := s.users[chatID]
	if !exists || state.PhoneNum == "" {
		s.sendMessage(chatID, "Please enter your phone number first using /start.")
		return
	}

	todayMenu, err := nutritionpro.GetMenuForToday(state.PhoneNum)
	if err != nil {
		log.Printf("Error getting menu for phone number %s: %v", state.PhoneNum, err)
		s.sendMessage(chatID, "Error retrieving today's menu. Please try again later.")
		return
	}
	s.sendMessage(chatID, todayMenu)
}

func (s *BotService) handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	if message.IsCommand() {
		s.handleCommand(message)
		return
	}

	state, exists := s.users[chatID]
	if !exists {
		s.initializeNewUser(chatID)
		return
	}

	if state.Stage == 1 {
		s.handlePhoneNumberInput(chatID, strings.TrimSpace(message.Text))
		return
	}

	s.handleHelp(chatID)
}

func (s *BotService) handleCommand(message *tgbotapi.Message) {
	command := message.Command()
	chatID := message.Chat.ID

	switch command {
	case "start":
		s.handleStart(chatID)
	case "help":
		s.handleHelp(chatID)
	case "unsubscribe":
		s.handleUnsubscribe(chatID)
	case "menu":
		s.handleMenuForToday(chatID)
	default:
		s.sendMessage(chatID, "Unknown command. Type /help for available commands.")
	}
}

func (s *BotService) handlePhoneNumberInput(chatID int64, phoneNumber string) {
	state := s.users[chatID]
	if isValidPhoneNumber(phoneNumber) {
		phoneNumberWithoutPrefix := strings.Replace(phoneNumber, "+420", "", 1)
		state.PhoneNum = phoneNumberWithoutPrefix
		state.Stage = 2

		// register for notifications
		if err := s.notifier.RegisterUser(chatID, phoneNumberWithoutPrefix); err != nil {
			log.Printf("Failed to register notifications for user %d: %v", chatID, err)
			s.sendMessage(chatID, "✅ Phone number saved, but notifications couldn't be enabled.")
			return
		}

		msg := strings.Join([]string{
			"✅ Phone number successfully saved: " + phoneNumber,
			"🔔 You will receive notifications for:",
			"• Breakfast at 7:00",
			"• Lunch at 11:00",
			"• Dinner at 17:00",
		}, "\n")

		s.sendMessage(chatID, msg)
		s.handleHelp(chatID)
	} else {
		errorText := strings.Join([]string{
			"❌ Invalid phone number format.",
			"Please enter a valid Czech phone number:",
			"• Must start with +420",
			"• Followed by 9 digits",
			"• Example: +420123456789",
		}, "\n")
		s.sendMessage(chatID, errorText)
	}
}

func (s *BotService) handleStart(chatID int64) {
	delete(s.users, chatID)
	s.initializeNewUser(chatID)
}

func (s *BotService) handleHelp(chatID int64) {
	helpText := strings.Join([]string{
		"Available commands:",
		"/help - Show this help message",
		"/start - Start over and enter a new phone number",
		"/unsubscribe - Remove your phone number from the system",
		"/menu - Get today's menu for your phone number",
	}, "\n")
	s.sendMessage(chatID, helpText)
}

func (s *BotService) handleUnsubscribe(chatID int64) {
	if _, exists := s.users[chatID]; exists {
		s.notifier.UnregisterUser(chatID)
		delete(s.users, chatID)
		s.sendMessage(chatID, "Your phone number has been removed and notifications disabled.")
	} else {
		s.sendMessage(chatID, "You don't have any phone number saved in the system.")
	}
}

func (s *BotService) initializeNewUser(chatID int64) {
	state := &UserState{
		Stage:  1,
		ChatID: chatID,
	}
	s.users[chatID] = state
	welcomeText := strings.Join([]string{
		"Welcome! This is simple bot, for checking your ordered menu from nutritionpro.cz:",
		"You can get today's menu by entering your phone number.",
		"Also, you'll receive notifications about your menu every day (in the morning), and during the day, for each meal individually. You can unsubscribe at any time.",
		"To get started, please enter your Czech phone number in the following format:",
		"• Must start with +420",
		"• Followed by 9 digits",
		"• Example: +420123456789",
	}, "\n")
	s.sendMessage(chatID, welcomeText)
}

func (s *BotService) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := s.api.Send(msg); err != nil {
		log.Printf("Error sending message to %d: %v", chatID, err)
	}
}
