package models

import (
	"time"

	"gorm.io/gorm"
)

type TicketStatus string

const (
	StatusPending    TicketStatus = "Pending"
	StatusInProgress TicketStatus = "In Progress"
	StatusResolved   TicketStatus = "Resolved"
)

type Conversation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"size:50;index" json:"user_id"`
	Message   string         `gorm:"type:text" json:"message"`
	Sentiment string         `gorm:"size:20" json:"sentiment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Ticket struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ConversationID uint           `gorm:"index" json:"conversation_id"`
	Description    string         `gorm:"size:255" json:"description"`
	Status         TicketStatus   `gorm:"type:varchar(20);default:'Pending'" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
