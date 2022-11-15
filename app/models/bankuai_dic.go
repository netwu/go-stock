package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

// Symbols ...
type BankuaiDic struct {
	orm.Model
	Id        int       `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null;index:idx_name,unique"`
	Count     int32     `gorm:"type:int;not null"`
	CreatedAt time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP"`
}
