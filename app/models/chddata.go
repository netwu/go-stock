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

func (m *Chddata) Xg() {
	// var results []map[string]interface{}
	type DataResults struct {
		Date string
	}
	type SymbolResults struct {
		Code string
		Name string
	}
	var dataResults []DataResults
	var symbolResults []SymbolResults
	facades.Orm.Query().Table("chddata").Select("distinct(date) as date").Order("date desc").Scan(&dataResults)
	if len(dataResults) > 0 {
		facades.Orm.Query().Table("symbols").Join("left join chddata c1 on symbols.code=c1.code left join chddata c2 on symbols.code=c2.code").Where("c1.date", dataResults[1].Date).Where("c2.date", dataResults[0].Date).Where("c2.`voturnover`/c1.voturnover >2 and c1.name  not like '%ST%' and c2.pchg<10 ").Select("symbols.code,symbols.name").Order("c1.turnover desc").Scan(&symbolResults)
		var codes []string
		for _, v := range symbolResults {
			codes = append(codes, v.Code)
		}

		facades.Orm.Query().Table("symbols as s").Join("join chddata c1 on s.code=c1.code join chddata c2 on s.code=c2.code join chddata c3 on s.code=c3.code ").Where("c1.date", dataResults[40].Date).Where("c2.date", dataResults[20].Date).Where("c3.date", dataResults[0].Date).Where("c1.tclose<c2.tclose and c2.tclose<c3.tclose").Where("s.code in ?", codes).Select("s.code,s.name").Order("c1.turnover desc").Scan(&symbolResults)
		for _, v := range symbolResults {
			fmt.Println(v.Code, v.Name)
		}

	}
}
