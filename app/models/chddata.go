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
	Id         int     `gorm:"primaryKey"`
	Date       string  `gorm:"type:date;not null;uniqueIndex:idx_code_date"`
	Market     string  `gorm:"type:varchar(255);not null"`
	Code       string  `gorm:"type:varchar(255);not null;uniqueIndex:idx_code_date"`
	Name       string  `gorm:"type:varchar(255);not null"`
	Tclose     float64 `gorm:"type:float;comment:收盘价;not null"`
	High       float64 `gorm:"type:float;comment:最高价;not null"`
	Low        float64 `gorm:"type:float;comment:最低价;not null"`
	Topen      float64 `gorm:"type:float;comment:开盘价;not null"`
	Chg        float64 `gorm:"type:float;comment:涨跌额;not null"`
	Pchg       float64 `gorm:"type:float;comment:涨跌幅;not null"`
	Turnover   float64 `gorm:"type:float;comment:换手率;not null"`
	Voturnover float64 `gorm:"type:float;comment:成交量;not null"`
	Vaturnover float64 `gorm:"type:float;comment:成交金额;not null"`
	Tcap       float64 `gorm:"type:float;comment:总市值;not null"`
	Mcap       float64 `gorm:"type:float;comment:流通市值;not null"`
	Amplitude  float64 `gorm:"type:float;comment:振幅;not null"`

	CreatedAt time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP on update current_timestamp;"`
}

func (m *Chddata) Store(datas []Chddata) (err error) {
	for _, data := range datas {
		sql := "INSERT INTO chddata (date,market,code,name,tclose,high,low,topen,chg,pchg,turnover,voturnover,vaturnover,tcap,mcap,amplitude) VALUES( '%s','%s','%s','%s',%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f ) ON DUPLICATE KEY UPDATE tclose =%f,high=%f,low =%f,topen =%f,chg =%f,pchg=%f,turnover =%f,voturnover=%f,vaturnover=%f,tcap=%f,mcap=%f,amplitude=%f"
		s := fmt.Sprintf(sql, data.Date, data.Market, data.Code, data.Name, data.Tclose, data.High, data.Low, data.Topen, data.Chg, data.Pchg, data.Turnover, data.Voturnover, data.Vaturnover, data.Tcap, data.Mcap, data.Amplitude, data.Tclose, data.High, data.Low, data.Topen, data.Chg, data.Pchg, data.Turnover, data.Voturnover, data.Vaturnover, data.Tcap, data.Mcap, data.Amplitude)
		facades.Log.Info("bankuaidic sql", s)
		facades.Orm.Query().Exec(s)
	}

	return nil
}

func (m *Chddata) UpdateChddataMonth() {
	sql := "INSERT INTO chddata_months(code,name,month,avg_price) SELECT sb.code,sb.name,sb.month,sb.avg_tclose FROM (SELECT code,name,left(date,7)as month,avg(tclose) as avg_tclose FROM chddata GROUP BY code,month) as sb ON DUPLICATE KEY UPDATE avg_price=sb.avg_tclose;"
	facades.Orm.Query().Exec(sql)
}
