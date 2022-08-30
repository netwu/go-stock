package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type Guzhis struct {
	orm.Model
	Id        int     `gorm:"primaryKey"`
	Symbol    string  `gorm:"type:varchar(255);not null"`
	Code      string  `gorm:"type:varchar(255);not null;index:idx_code,unique"`
	Name      string  `gorm:"type:varchar(255);not null"`
	Date      string  `gorm:"type:date;not null;index:idx_date"`
	Hangye    string  `gorm:"type:varchar(255);not null"`
	Zhuying   string  `gorm:"type:varchar(255);not null"`
	Low_pe    float64 `gorm:"type:float;not null"`
	Std_pe    float64 `gorm:"type:float;not null"`
	High_pe   float64 `gorm:"type:float;not null"`
	Warm_pe   float64 `gorm:"type:float;not null"`
	Normal_pe float64 `gorm:"type:float;not null"`

	Low_price    float64   `gorm:"type:float;not null"`
	Std_price    float64   `gorm:"type:float;not null"`
	High_price   float64   `gorm:"type:float;not null"`
	Warm_price   float64   `gorm:"type:float;not null"`
	Normal_price float64   `gorm:"type:float;not null"`
	Cbs_score    float64   `gorm:"type:float;not null"`
	Cbs_rank     float64   `gorm:"type:float;not null"`
	Roe          float64   `gorm:"type:float;not null"`
	CreatedAt    time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP"`
}
