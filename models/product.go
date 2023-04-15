package models

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Product struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Title     string     `gorm:"not null" json:"title" valid:"required~Title is required"`
	Caption   string     `json:"caption" valid:"required~Caption is required"`
	Price     int64      `json:"price" valid:"required~Price is required"`
	Photo_url string     `json:"photo_url" valid:"required~Photo URL is required"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (u *Product) BeforeCreate(tx *gorm.DB) (err error) {

	_, errCreate := govalidator.ValidateStruct(u)
	if errCreate != nil {
		err = errCreate
		return
	}
	return
}

func (u *Product) BeforeUpdate(tx *gorm.DB) (err error) {

	_, errCreate := govalidator.ValidateStruct(u)
	if errCreate != nil {
		err = errCreate
		return
	}

	return
}

func (u *Product) ToString() string {
	return fmt.Sprintf("id : %d\nTitle : %s\nCaption : %s\nPhoto_url : %s\nCreated_at : %s\nUpdated_at : %s", u.ID, u.Title, u.Caption, u.Photo_url, u.CreatedAt, u.UpdatedAt)
}
