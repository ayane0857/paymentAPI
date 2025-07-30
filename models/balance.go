package models

import "time"

type Balance struct {
  ID               uint           `gorm:"primary_key"`
  Balance          int            `gorm:"size:255"`
  CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
  UpdatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
}

type BalanceResponse struct {
  Balance    float64   `json:"balance"`     // `json:"balance"`でJSONキー名を指定
  TotalMoney float64   `json:"total_money"` // こちらも同様
  CreatedAt  time.Time `json:"created_at"`
  UpdatedAt  time.Time `json:"updated_at"`
}