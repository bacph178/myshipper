package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	Name        string       `gorm:"size:280;not null"`
	Slug        string       `gorm:"unique_index;not null"`
	Price       int          `gorm:"not null"`
	Stock       int          `gorm:"not null"`
	Comments    []Comment    `gorm:"foreignKey:ProductId"`
	Tags        []Tag        `gorm:"many2many:products_tags;"`
	Images      []FileUpload `gorm:"foreignKey:ProductId"`
	Categories  []Category   `gorm:"many2many:products_categories;"`
	Description string       `gorm:"not null"`
}

func (product *Product) BeforeSave() (err error) {
	product.Slug = slug.Make(product.Name)
	return
}
