package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type School struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Address   string         `json:"address" gorm:"not null"`
	Type      string         `json:"type" gorm:"not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Reviews   []Review       `json:"reviews,omitempty" gorm:"foreignKey:SchoolID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Review struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Rating    float32        `json:"rating" gorm:"not null"`
	Comment   string         `json:"comment"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	SchoolID  uint           `json:"school_id" gorm:"not null"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	School    School         `json:"-" gorm:"foreignKey:SchoolID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
