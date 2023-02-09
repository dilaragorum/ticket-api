package ticket

import (
	"gorm.io/gorm"
)

type Ticket struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"not null;unique" json:"name"`
	Desc       string `gorm:"not null" json:"desc"`
	Allocation int    `gorm:"not null;check:allocation>0" json:"allocation"`
	gorm.Model
}

type Purchase struct {
	ID       int `gorm:"primaryKey"`
	UserID   string
	TicketID int `gorm:"not null"`
	Quantity int `gorm:"not null;check:quantity>0"`
	gorm.Model
}

func (Purchase) TableName() string {
	return "tickets_purchases"
}
