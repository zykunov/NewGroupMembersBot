# Бот для telegram на Golang
## позволяет отслеживать новых пользователей в группах VK

Для манипуляций с БД (PostgreSQL) используется GORM

Для взаимодействия с API Telegram - библиотека https://github.com/go-telegram-bot-api/telegram-bot-api

Для доступа к VK API - просто запрос через http.Get, нам нужен всего один метод, поэтому не стал усложнять и множить зависимости.

### Директории
**cmd/vkapibot/main.go** - точка входа

**configs/** - кофиги БД и токены VK и TG

**internal/app/botcore.go** - логика работы телеграмной части приложения

**internal/app/vkcore.go** - логика работы с API вконтакте.

**models/** - модели GORM

**pkg/keyboardgenerator/keyboardgenerator.go** - генератор клавиатуры в TG, мб пригодится в других проектах.

**storage/** - репозитории описывающие функции для моделей user и group, а так же инициализация соединения с БД и миграции(storage.go).
