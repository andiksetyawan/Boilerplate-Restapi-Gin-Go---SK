package models

import "github.com/jinzhu/gorm"

type Article struct {
    gorm.Model
    Title string
    Slug string `gorm:"unique_index"`
    Desc string `sql:"type:text;"`
    Tag string
    UserID  uint
}