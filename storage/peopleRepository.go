package storage

import (
	"log"
	"strconv"
	"strings"

	"github.com/zykunov/courseGoFirst/vkApiBot/models"
)

// добавление всех пользователей. Принимает слайс из моделей users и добавляет в бд
func UserAdd(userSlice []models.User) (msg string) {

	var stBuilder strings.Builder

	for _, user := range userSlice {
		string := "https://vk.com/id" + strconv.Itoa(user.VkId) + "\n"
		stBuilder.WriteString(string)

		GetDB().Create(&user)
	}

	if stBuilder.String() == "" {
		return "Нет новых пользователей"
	}

	return "Группа добавлена, подписчики загружены 👌"
}

// селект всех пользователей
func GetAllUsers(userId int64, vkgroup_id string) (usersFromDB []*models.User) {

	result := make([]*models.User, 0)

	err := GetDB().Table("users").Select("vk_id").Where("user_id = ? AND vkgroup_id = ?", int(userId), vkgroup_id).Find(&result).Error
	if err != nil {
		log.Println(err)
	}

	return result
}

// Удаление пользователей
func DeleteUsersByGroupId(userId int64, groupId string) (msg string) {

	result := make([]*models.User, 0)

	err := GetDB().Table("users").Where("user_id = ? AND vkgroup_id LIKE ? ", int(userId), groupId).Find(&result).Error
	if err != nil {
		log.Println(err)
	}

	err = GetDB().Unscoped().Delete(&result).Error //Безвозвратное удаление
	if err != nil {
		return "Не получилось удалить пользователей группы"
	}

	return "Группа " + groupId + " удалена!"
}
