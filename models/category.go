package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Category struct {
	gorm.Model
	Name        string       `gorm:"not null"`
	Description string       `gorm:"default:null"`
	Slug        string       `gorm:"unique_index"`
	Products    []Product    `gorm:"many2many:products_categories;"`
	Images      []FileUpload `gorm:"foreignKey:CategoryId"`
	IsNewRecord bool         `gorm:"-;default:false"`
}

func (category *Category) BeforeSave() (err error) {
	category.Slug = slug.Make(category.Name)
	return
}
