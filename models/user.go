package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id        uint   `gorm:"primarykey"`
	VkgroupId string //id группы
	VkId      int    //id пользователя VK
	UserID    int    // uid из tg
	Sex       int    // пол
	IsNew     bool
}
