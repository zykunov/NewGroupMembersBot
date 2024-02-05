package app

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewKeyBoard(buttons string) tgbotapi.ReplyKeyboardMarkup {

	buttonSlice := string2Slice(buttons)

	var keyboard [][]tgbotapi.KeyboardButton
	var button tgbotapi.KeyboardButton
	var rows []tgbotapi.KeyboardButton

	for _, value := range buttonSlice {
		var rows []tgbotapi.KeyboardButton

		fmt.Println(value)

		button = tgbotapi.KeyboardButton{
			Text:            value,
			RequestLocation: false,
		}
		rows = append(rows, button)

		keyboard = append(keyboard, rows)
	}

	button = tgbotapi.KeyboardButton{
		Text: "⬅️ Назад",
	}
	rows = append(rows, button)
	keyboard = append(keyboard, rows)
	return tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard:       keyboard,
	}
}

func string2Slice(buttonsString string) []string {

	// buttonsString = strings.ReplaceAll(buttonsString, " ", "")
	buttonSlice := strings.Split(buttonsString, "\n")

	return buttonSlice
}
