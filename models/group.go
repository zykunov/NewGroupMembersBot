package models

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Id        uint `gorm:"primarykey"`
	VkgroupId string
	UserID    int
}
