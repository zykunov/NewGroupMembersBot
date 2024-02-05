package main

import (
	"github.com/zykunov/courseGoFirst/vkApiBot/internal/app"
	"github.com/zykunov/courseGoFirst/vkApiBot/storage"
)

func init() {
	storage.GetDB()
}

func main() {

	app.BotStart()

}
