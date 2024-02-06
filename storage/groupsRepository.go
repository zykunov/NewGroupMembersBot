package storage

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strings"

	"github.com/zykunov/courseGoFirst/vkApiBot/models"
)

// Вывести все группы
func GetAllGroups(userId int64) (groupsList string) {
	result := make([]*models.Group, 0)
	err := GetDB().Table("groups").Select("vkgroup_id", "user_id").Where("user_id = ?", int(userId)).Find(&result).Error
	if err != nil {
		fmt.Println(err)
	}

	var stBuilder strings.Builder
	for _, value := range result {
		// string := "https://vk.com/" + value.VkgroupId + "\n"
		string := value.VkgroupId + "\n"

		stBuilder.WriteString(string)
	}

	if stBuilder.String() == "" {
		return "Вы не добавили ни одной группы для отслеживания"
	}
	return stBuilder.String()
}

// Добавление групп
func GroupAdd(userId int64, groupLink string) (msg string, groupPath string, err error) {

	if reflect.TypeOf(groupLink).String() != "string" {
		return "введите строковое значение", "", err
	}

	_, err = url.ParseRequestURI(groupLink)
	if err != nil {
		return "Ссылка не распознана", "", err
	}

	parsedURL, err := url.Parse(groupLink)
	if err != nil {
		log.Println("Can't parse link from user")
		return "Не удалось распознать ссылку", "", err
	}
	path := strings.ReplaceAll(parsedURL.Path, `/`, "")

	log.Println("Trying write to DB", groupLink)

	group := &models.Group{
		UserID:    int(userId),
		VkgroupId: path,
	}

	GetDB().Create(group)

	return "Группа " + groupLink + " добавлена", path, nil
}

// Удаление групп
func DeleteGroupById(userId int64, groupLink string) (msg string) {

	if reflect.TypeOf(groupLink).String() != "string" {
		return "введите строковое значение"
	}

	_, err := url.ParseRequestURI(groupLink)
	if err != nil {
		return "Ссылка не распознана"
	}

	parsedURL, err := url.Parse(groupLink)
	if err != nil {
		log.Println("(delete group method) Can't parse link from user")
		return "Ссылка не распознана"
	}
	path := strings.ReplaceAll(parsedURL.Path, `/`, "")

	result := make([]*models.Group, 0)
	err = GetDB().Table("groups").Where("user_id = ? AND vkgroup_id LIKE ? ", int(userId), path).Find(&result).Error
	if err != nil {
		fmt.Println(err)
	}

	err = GetDB().Delete(&result).Error // мягкое удаление
	if err != nil {
		return "Группа не найдена"
	}

	return "Группа " + groupLink + " удалена"
}

// Получаем группу по id
func CheckGroupById(groupId string, userId int64) (checkGroup bool) {
	var count int64
	var check bool
	err := GetDB().Table("groups").Select("vkgroup_id", "user_id").Where("user_id = ? AND vkgroup_id LIKE ?", int(userId), groupId).Count(&count)
	if err != nil {
		fmt.Println(err)
	}

	if count != 0 {
		check = true
		return check
	}

	check = false
	return check
}
