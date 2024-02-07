package app

import (
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	keyBoardGenerator "github.com/zykunov/courseGoFirst/vkApiBot/pkg/keyboardgenerator"
	"github.com/zykunov/courseGoFirst/vkApiBot/storage"
)

var (
	linkWait       bool
	delete         bool
	groupMode      bool
	previusMessage string
)

func BotStart() {
	botToken := os.Getenv("bot_token")
	if botToken == "" {
		log.Fatalf("set token for your bot!")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("error while starting bot", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	var numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Управление группами"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✅Добавить группу для отслеживания"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❓Как этим пользоваться?"),
		),
	)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%v:%s] %s", update.Message.From.ID, update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			if groupMode {
				if update.Message.Text == "⬅️ Назад" {
					msg.ReplyMarkup = numericKeyboard
					groupMode = false
				}
				groupExist := storage.CheckGroupById(update.Message.Text, update.Message.Chat.ID)
				if groupExist {
					msg.ReplyMarkup = NewKeyBoard("🆕Смотреть новых людей\n❌Удалить группу")
					msg.Text = "Выберите действия с группой"
					groupMode = false
					log.Println("disable group mode")
					previusMessage = update.Message.Text
				} else {
					msg.Text = "Такой группы нет"
				}
			} else {
				if linkWait {
					if delete {
						msg.DisableWebPagePreview = true
						msg.Text = storage.DeleteGroupById(update.Message.Chat.ID, update.Message.Text)
						linkWait, delete = false, false
					} else {
						var groupPath string
						msg.DisableWebPagePreview = true

						msg.Text, groupPath, err = storage.GroupAdd(update.Message.Chat.ID, update.Message.Text) //1 добавляем группу
						if err != nil {
							break
						}

						log.Println("geting group members from vk api")
						users := GetGroupMembers(groupPath, int(update.Message.Chat.ID)) //2 сразу получаем её пользоватлей

						log.Println("write users from vk api2DB")
						storage.UserAdd(users) //3 и добавляем пользователей в БД

						linkWait = false
					}
				} else {
					switch update.Message.Text {
					case "/start":
						msg.ReplyMarkup = numericKeyboard
					case "Управление группами":
						groupMode = true
						log.Println("group mode activated")
						groups := storage.GetAllGroups(update.Message.Chat.ID)
						msg.ReplyMarkup = keyBoardGenerator.NewKeyBoard(groups)
					case "Посмотреть список групп":
						msg.DisableWebPagePreview = true
						msg.Text = storage.GetAllGroups(update.Message.Chat.ID)
					case "✅Добавить группу для отслеживания":
						msg.Text = "Вставьте ссылку на группу"
						linkWait = true
					case "⬅️ Назад":
						msg.ReplyMarkup = numericKeyboard
						groupMode = false
					case "❌Удалить группу":
						msg.Text = storage.DeleteGroupById(update.Message.Chat.ID, "https://vk.com/"+previusMessage)
						msg.Text = storage.DeleteUsersByGroupId(update.Message.Chat.ID, previusMessage)
						log.Printf("group name:", previusMessage)
						linkWait = true
						delete = true
						previusMessage = ""
					case "🆕Смотреть новых людей":
						msg.DisableWebPagePreview = true
						usersFromDB := storage.GetAllUsers(update.Message.Chat.ID, previusMessage)
						usersFromVK := GetGroupMembers(previusMessage, int(update.Message.Chat.ID))

						var stBuilder strings.Builder

						for _, valueApi := range usersFromVK {

							var find bool

							for _, valueDB := range usersFromDB {

								if valueApi.VkId == valueDB.VkId {
									/* Сделать добавление польз в Базу */
									find = true
									break
								}
							}

							if !find {
								log.Println("new user finded!")
								string := "https://vk.com/id" + strconv.Itoa(valueApi.VkId) + "\n"
								stBuilder.WriteString(string)
							}
						}
						if stBuilder.String() == "" {
							msg.Text = "Новых пользователей нет"
						} else {
							stBuilder.WriteString("⬆️новые пользователи⬆️")
							msg.Text = stBuilder.String()
						}
					case "❓Как этим пользоваться?":
						msg.Text = "Бот может отслеживать новых пользователей в группах вк. (Максимум 5 групп)\n 👉 Новые пользователи считаются со времени последнего просмотра"
					default:
						msg.Text = "Не знаю таких команд, можете попробовать /start"
					}
				}
			}

			if _, err := bot.Send(msg); err != nil {
				log.Printf("can't send message")
			}

		}
	}
}
