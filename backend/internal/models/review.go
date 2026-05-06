package models

import (
	"database/sql"
	"time"
)

type Review struct {
	ID                 uint64         `gorm:"primaryKey" json:"id"`
	ProductID          uint64         `gorm:"not null;index" json:"product_id"`
	UserID             uint64         `gorm:"not null;index" json:"user_id"`
	OrderID            *uint64        `json:"order_id,omitempty"`
	Rating             int            `gorm:"not null;check:rating BETWEEN 1 AND 5" json:"rating"`
	Title              *string        `gorm:"type:varchar(255)" json:"title,omitempty"`
	Comment            *string        `gorm:"type:text" json:"comment,omitempty"`
	Images             sql.NullJSON   `gorm:"type:json" json:"images,omitempty"`
	IsVerifiedPurchase bool           `gorm:"default:false" json:"is_verified_purchase"`
	IsApproved         bool           `gorm:"default:false" json:"is_approved"`
	HelpfulCount       int            `gorm:"default:0" json:"helpful_count"`
	RepliedBy          *uint64        `json:"replied_by,omitempty"`
	ReplyText          *string        `gorm:"type:text" json:"reply_text,omitempty"`
	RepliedAt          *time.Time     `json:"replied_at,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`

	// Relationships
	Product Product `json:"product,omitempty"`
	User    User    `json:"user,omitempty"`
	Order   *Order  `json:"order,omitempty"`
	Helpful []ReviewHelpful `json:"-"`
}

func (Review) TableName() string {
	return "reviews"
}

type ReviewHelpful struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	ReviewID   uint64    `gorm:"not null;index" json:"review_id"`
	UserID     uint64    `gorm:"not null;index" json:"user_id"`
	IsHelpful  bool      `gorm:"default:true" json:"is_helpful"`
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	Review Review `json:"-"`
	User   User   `json:"-"`
}

func (ReviewHelpful) TableName() string {
	return "review_helpfuls"
}