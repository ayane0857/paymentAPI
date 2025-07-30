package models

import "time"

type Payment struct {
    ID            uint      `gorm:"primary_key"`
    Money         int       `gorm:"size:255"`
    PaymentMethod string    `gorm:"size:255"`
    CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
    UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type PaymentRequest struct {
    Money         *int    `json:"money"`
    PaymentMethod *string `json:"paymentMethod"`
}