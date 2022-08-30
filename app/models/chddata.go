package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/facades"
	"gorm.io/gorm/clause"
)

// Symbols ...
type Chddata struct {
	orm.Model
	Id         int       `gorm:"primaryKey"`
	Date       string    `gorm:"type:date;not null;uniqueIndex:idx_code_date"`
	Symbol     string    `gorm:"type:varchar(255);not null"`
	Code       string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_code_date"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Tclose     float64   `gorm:"type:float;comment:收盘价;not null"`
	High       float64   `gorm:"type:float;comment:最高价;not null"`
	Low        float64   `gorm:"type:float;comment:最低价;not null"`
	Topen      float64   `gorm:"type:float;comment:开盘价;not null"`
	Chg        float64   `gorm:"type:float;comment:涨跌额;not null"`
	Pchg       float64   `gorm:"type:float;comment:涨跌幅;not null"`
	Turnover   float64   `gorm:"type:float;comment:换手率;not null"`
	Voturnover float64   `gorm:"type:float;comment:成交量;not null"`
	Vaturnover float64   `gorm:"type:float;comment:成交金额;not null"`
	Tcap       float64   `gorm:"type:float;comment:总市值;not null"`
	Mcap       float64   `gorm:"type:float;comment:流通市值;not null"`
	CreatedAt  time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP"`
}

func (m *Chddata) Store(a *[]Chddata) (err error) {
	facades.DB.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(a)
	// ret := m.Conn.Create(a)
	return nil
}
