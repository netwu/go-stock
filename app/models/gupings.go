package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type Gupings struct {
	orm.Model
	Id        int       `gorm:"primaryKey"`
	Date      string    `gorm:"type:date;not null;index:idx_date"`
	Symbol    string    `gorm:"type:varchar(255);not null"`
	Code      string    `gorm:"type:varchar(255);not null;index:idx_code"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Guping    string    `gorm:"type:varchar(512);not null"`
	Short     string    `gorm:"type:varchar(512);not null"`
	Mid       string    `gorm:"type:varchar(512);not null"`
	Long      string    `gorm:"type:varchar(512);not null"`
	Score     float64   `gorm:"type:float;not null"`
	Opt       string    `gorm:"type:varchar(8);not null"`
	CreatedAt time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP"`
}
