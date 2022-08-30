package models

import "github.com/goravel/framework/database/orm"

// Symbols ...
type Bankuais struct {
	orm.Model
	Id      int    `gorm:"primaryKey"`
	Symbol  string `gorm:"type:varchar(255);not null"`
	Code    string `gorm:"type:varchar(255);not null;index:idx_code,unique"`
	Name    string `gorm:"type:varchar(255);not null"`
	Bankuai string `gorm:"type:varchar(512);not null"`
	Zhuying string `gorm:"type:varchar(512);not null"`
}
