package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Tag struct {
	gorm.Model
	Name        string       `gorm:"not null"`
	Description string       `gorm:"default:null"`
	Slug        string       `gorm:"unique_index"`
	Products    []Product    `gorm:"many2many:products_tags;"`
	Images      []FileUpload `gorm:"foreignKey:TagId"`
	IsNewRecord bool         `gorm:"-;default:false"`
}

func (tag *Tag) BeforeSave() (err error) {
	tag.Slug = slug.Make(tag.Name)
	return
}
