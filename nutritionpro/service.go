package nutritionpro

import (
	"fmt"
	"strings"
	"time"
)

var numberMenuMap = make(map[string]*MenuResponse)

func GetMenu(phoneNumber string) (*MenuResponse, error) {
	if menu, exists := numberMenuMap[phoneNumber]; exists {
		return menu, nil
	}

	client, err := NewClientWithAuth(phoneNumber)
	if err != nil {
		return nil, err
	}

	menu, err := client.GetMenu()
	if err != nil {
		return nil, err
	}

	//mockClient := NewMockClient()
	//menu, _ := mockClient.GetMenu()

	numberMenuMap[phoneNumber] = menu
	return menu, nil
}

func GetMenuForToday(phoneNumber string) (string, error) {
	menu, err := GetMenu(phoneNumber)
	if err != nil {
		return "", err
	}

	today := time.Now().Truncate(24 * time.Hour).Unix()
	var todayMenu *Day
	for _, day := range menu.Days {
		if day.Timestamp == today {
			todayMenu = &day
			break
		}
	}

	if todayMenu == nil {
		return "🌱 No menu available for today, enjoy your day", nil
	}

	return FormatDayMenu(*todayMenu), nil
}

func ClearMenu(phoneNumber string) {
	delete(numberMenuMap, phoneNumber)
}

func FormatDayMenu(day Day) string {
	date := time.Unix(day.Timestamp, 0)

	var sb strings.Builder

	// Write header with date
	sb.WriteString(fmt.Sprintf("📅 Menu for %s, %s\n\n",
		date.Format("Monday"),
		date.Format("02.01.2006"),
	))

	// Sort dishes by meal time (0 = breakfast, 2 = lunch, 4 = dinner)
	mealNames := map[int]string{
		0: "🌅 Breakfast",
		2: "🍽️ Lunch",
		4: "🌙 Dinner",
	}

	// Write dishes grouped by meal time
	for _, dish := range day.Dishes {
		mealName := mealNames[dish.Meal]
		sb.WriteString(fmt.Sprintf("%s:\n", mealName))
		sb.WriteString(fmt.Sprintf("🍳 %s\n", dish.Title))

		if dish.Description != "" {
			sb.WriteString(fmt.Sprintf("ℹ️ %s\n", dish.Description))
		}

		sb.WriteString(fmt.Sprintf("⚖️ Weight: %dg", dish.Weight))
		if dish.IsHot {
			sb.WriteString(" (♨️ Needs heating)")
		}
		sb.WriteString("\n")

		sb.WriteString(fmt.Sprintf("📊 Nutrients:\n"))
		sb.WriteString(fmt.Sprintf("   • Calories: %.0f kcal\n", dish.Nutrients.Kcal))
		sb.WriteString(fmt.Sprintf("   • Protein: %.1fg\n", dish.Nutrients.Prot))
		sb.WriteString(fmt.Sprintf("   • Fat: %.1fg\n", dish.Nutrients.Fat))
		sb.WriteString(fmt.Sprintf("   • Carbs: %.1fg\n", dish.Nutrients.Carb))
		sb.WriteString("\n")
	}

	sb.WriteString("📈 Daily Totals:\n")
	sb.WriteString(fmt.Sprintf("   • Total Calories: %.0f kcal\n", day.Nutrients.Kcal))
	sb.WriteString(fmt.Sprintf("   • Total Protein: %.1fg\n", day.Nutrients.Prot))
	sb.WriteString(fmt.Sprintf("   • Total Fat: %.1fg\n", day.Nutrients.Fat))
	sb.WriteString(fmt.Sprintf("   • Total Carbs: %.1fg\n", day.Nutrients.Carb))

	return strings.TrimSpace(sb.String())
}
