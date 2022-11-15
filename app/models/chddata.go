package models

import (
	"fmt"
	"time"

	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
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

func (m *Chddata) Store(datas []Chddata) (err error) {
	for _, data := range datas {
		sql := "INSERT INTO chddata (date,symbol,code,name,tclose,high,low,topen,chg,pchg,turnover,voturnover,vaturnover,tcap,mcap) VALUES( '%s','%s','%s','%s',%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f ) ON DUPLICATE KEY UPDATE tclose =%f,high=%f,low =%f,topen =%f,chg =%f,pchg=%f,turnover =%f,voturnover=%f,vaturnover=%f,tcap=%f,mcap=%f"
		s := fmt.Sprintf(sql, data.Date, data.Symbol, data.Code, data.Name, data.Tclose, data.High, data.Low, data.Topen, data.Chg, data.Pchg, data.Turnover, data.Voturnover, data.Vaturnover, data.Tcap, data.Mcap, data.Tclose, data.High, data.Low, data.Topen, data.Chg, data.Pchg, data.Turnover, data.Voturnover, data.Vaturnover, data.Tcap, data.Mcap)
		facades.Log.Info("bankuaidic sql", s)
		facades.Orm.Query().Exec(s)
	}

	return nil
}
