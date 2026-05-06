package models

import (
	"time"
)

type Address struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	UserID         uint64     `gorm:"not null;index" json:"user_id"`
	Label          string     `gorm:"type:varchar(50);default:'Rumah'" json:"label"`
	RecipientName  string     `gorm:"type:varchar(100);not null" json:"recipient_name"`
	Phone          string     `gorm:"type:varchar(20);not null" json:"phone"`
	Province       string     `gorm:"type:varchar(100);not null" json:"province"`
	ProvinceID     *string    `gorm:"type:varchar(10)" json:"province_id,omitempty"`
	City           string     `gorm:"type:varchar(100);not null" json:"city"`
	CityID         *string    `gorm:"type:varchar(10)" json:"city_id,omitempty"`
	District       *string    `gorm:"type:varchar(100)" json:"district,omitempty"`
	SubdistrictID  *string    `gorm:"type:varchar(10)" json:"subdistrict_id,omitempty"`
	PostalCode     string     `gorm:"type:varchar(10);not null" json:"postal_code"`
	AddressDetail  string     `gorm:"type:text;not null" json:"address_detail"`
	Landmark       *string    `gorm:"type:text" json:"landmark,omitempty"`
	Latitude       *float64   `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude      *float64   `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	IsDefault      bool       `gorm:"default:false" json:"is_default"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relationships
	User User `json:"user,omitempty"`
}

func (Address) TableName() string {
	return "addresses"
}