package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

// Symbols ...
type ChddataMonth struct {
	orm.Model
	Id        int       `gorm:"primaryKey"`
	Code      string    `gorm:"type:varchar(255);not null;index:idx_code_month,unique"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Month     string    `gorm:"type:varchar(255);not null;index:idx_code_month,unique"`
	AvgPrice  float64   `gorm:"type:float;not null"`
	CreatedAt time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP"`
}
