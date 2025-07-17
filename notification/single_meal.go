package notification

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"log"
	"nutritionpro-notify/nutritionpro"
	"sync"
	"time"
)

const (
	BreakfastHour = 7  // 7:00
	LunchHour     = 11 // 11:00
	DinnerHour    = 17 // 17:00
)

type MealNotifier struct {
	bot       *tgbotapi.BotAPI
	users     map[int64]*UserNotification
	scheduler *cron.Cron
	mutex     sync.RWMutex
}

type UserNotification struct {
	ChatID   int64
	PhoneNum string
	Client   *nutritionpro.Client
}

func NewMealNotifier(bot *tgbotapi.BotAPI) *MealNotifier {
	return &MealNotifier{
		bot:       bot,
		users:     make(map[int64]*UserNotification),
		scheduler: cron.New(cron.WithLocation(time.Local)),
	}
}

func (n *MealNotifier) RegisterUser(chatID int64, phoneNum string) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	client, err := nutritionpro.NewClientWithAuth(phoneNum)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	n.users[chatID] = &UserNotification{
		ChatID:   chatID,
		PhoneNum: phoneNum,
		Client:   client,
	}

	return nil
}

func (n *MealNotifier) UnregisterUser(chatID int64) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	delete(n.users, chatID)
}

func (n *MealNotifier) Start() error {
	// breakfast notification
	_, err := n.scheduler.AddFunc(fmt.Sprintf("0 %d * * *", BreakfastHour), func() {
		n.sendMealNotification(0) // 0 for breakfast
	})
	if err != nil {
		return fmt.Errorf("scheduling breakfast notification: %w", err)
	}

	// lunch notification
	_, err = n.scheduler.AddFunc(fmt.Sprintf("0 %d * * *", LunchHour), func() {
		n.sendMealNotification(2) // 2 for lunch
	})
	if err != nil {
		return fmt.Errorf("scheduling lunch notification: %w", err)
	}

	// dinner notification
	_, err = n.scheduler.AddFunc(fmt.Sprintf("0 %d * * *", DinnerHour), func() {
		n.sendMealNotification(4) // 4 for dinner
	})
	if err != nil {
		return fmt.Errorf("scheduling dinner notification: %w", err)
	}

	n.scheduler.Start()
	return nil
}

func (n *MealNotifier) Stop() {
	if n.scheduler != nil {
		n.scheduler.Stop()
	}
}

func (n *MealNotifier) sendMealNotification(mealType int) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	for _, user := range n.users {
		go func(u *UserNotification) {
			menu, err := u.Client.GetMenu()
			if err != nil {
				log.Printf("Error getting menu for user %d: %v", u.ChatID, err)
				return
			}

			// Find today's menu
			today := time.Now().Truncate(24 * time.Hour).Unix()
			var todayMenu *nutritionpro.Day
			for _, day := range menu.Days {
				if day.Timestamp == today {
					todayMenu = &day
					break
				}
			}

			if todayMenu == nil {
				log.Printf("No menu found for today for user %d", u.ChatID)
				return
			}

			// Find the specific meal
			var meal *nutritionpro.Dish
			for _, dish := range todayMenu.Dishes {
				if dish.Meal == mealType {
					meal = &dish
					break
				}
			}

			if meal == nil {
				log.Printf("No meal found for type %d for user %d", mealType, u.ChatID)
				return
			}

			// Send the notification
			message := formatMealMessage(meal)
			msg := tgbotapi.NewMessage(u.ChatID, message)
			if _, err := n.bot.Send(msg); err != nil {
				log.Printf("Error sending notification to user %d: %v", u.ChatID, err)
			}
		}(user)
	}
}

func formatMealMessage(meal *nutritionpro.Dish) string {
	var mealType string
	switch meal.Meal {
	case 0:
		mealType = "🌅 Breakfast"
	case 2:
		mealType = "🍽️ Lunch"
	case 4:
		mealType = "🌙 Dinner"
	}

	msg := fmt.Sprintf("%s time!\n\n", mealType)
	msg += fmt.Sprintf("🍳 %s\n\n", meal.Title)

	if meal.Description != "" {
		msg += fmt.Sprintf("ℹ️ %s\n\n", meal.Description)
	}

	msg += fmt.Sprintf("⚖️ Weight: %dg\n", meal.Weight)
	if meal.IsHot {
		msg += "♨️ Needs heating\n"
	}

	msg += "\n📊 Nutrients:\n"
	msg += fmt.Sprintf("   • Calories: %.0f kcal\n", meal.Nutrients.Kcal)
	msg += fmt.Sprintf("   • Protein: %.1fg\n", meal.Nutrients.Prot)
	msg += fmt.Sprintf("   • Fat: %.1fg\n", meal.Nutrients.Fat)
	msg += fmt.Sprintf("   • Carbs: %.1fg\n", meal.Nutrients.Carb)

	return msg
}
