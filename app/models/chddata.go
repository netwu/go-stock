package models

import (
	"fmt"
	"strings"
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
	if len(datas) == 0 {
		return nil
	}

	// 构建批量插入SQL
	var values []string
	for _, data := range datas {
		value := fmt.Sprintf("('%s','%s','%s','%s',%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f)",
			data.Date, data.Market, data.Code, data.Name,
			data.Tclose, data.High, data.Low, data.Topen,
			data.Chg, data.Pchg, data.Turnover,
			data.Voturnover, data.Vaturnover,
			data.Tcap, data.Mcap, data.Amplitude)
		values = append(values, value)
		// 删除今天的数据
		facades.Orm.Query().Table("chddata").
			Where("code = ?", data.Code).
			Where("date = ?", data.Date).
			Delete(&Chddata{})
	}

	sql := `INSERT INTO chddata (date,market,code,name,tclose,high,low,topen,chg,pchg,turnover,voturnover,vaturnover,tcap,mcap,amplitude) 
            VALUES %s`
	// sql := `INSERT INTO chddata (date,market,code,name,tclose,high,low,topen,chg,pchg,turnover,voturnover,vaturnover,tcap,mcap,amplitude)
	// VALUES %s
	// ON DUPLICATE KEY UPDATE
	// tclose=VALUES(tclose),
	// high=VALUES(high),
	// low=VALUES(low),
	// topen=VALUES(topen),
	// chg=VALUES(chg),
	// pchg=VALUES(pchg),
	// turnover=VALUES(turnover),
	// voturnover=VALUES(voturnover),
	// vaturnover=VALUES(vaturnover),
	// tcap=VALUES(tcap),
	// mcap=VALUES(mcap),
	// amplitude=VALUES(amplitude)`

	batchSQL := fmt.Sprintf(sql, strings.Join(values, ","))
	facades.Log.Info("执行批量更新SQL")

	facades.Orm.Query().Exec(batchSQL)

	return nil
}

func (m *Chddata) UpdateChddataMonth() {
	sql := "INSERT INTO chddata_months(code,name,month,avg_price) SELECT sb.code,sb.name,sb.month,sb.avg_tclose FROM (SELECT code,name,left(date,7)as month,avg(tclose) as avg_tclose FROM chddata GROUP BY code,month) as sb ON DUPLICATE KEY UPDATE avg_price=sb.avg_tclose;"
	facades.Orm.Query().Exec(sql)
}

func (m *Chddata) Xg(date string) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	// var results []map[string]interface{}
	type DataResults struct {
		Date string
	}
	type SymbolResults struct {
		Code   string
		Name   string
		Buy    float64
		Mktcap float64
	}
	var dataResults []DataResults
	var symbolResults []SymbolResults
	facades.Orm.Query().Table("chddata").Select("distinct(date) as date").Where("date<=?", date).Order("date desc").Scan(&dataResults)
	if len(dataResults) > 0 {
		facades.Orm.Query().Table("symbols s").Join("left join chddata c1 on s.code=c1.code left join chddata c2 on s.code=c2.code").Where("c1.date", dataResults[1].Date).Where("c2.date", dataResults[0].Date).Where("c2.`voturnover`/c1.voturnover >2 and c1.name not like '%ST%' and c1.code not like '%30%' and c1.code not like '%68%'  and c2.pchg<10 and c2.chg>0").Select("s.code,s.name").Order("c1.turnover desc").Scan(&symbolResults)
		var codes []string
		for _, v := range symbolResults {
			codes = append(codes, v.Code)
		}

		facades.Orm.Query().Table("symbols as s").Join("join chddata c1 on s.code=c1.code join chddata c2 on s.code=c2.code join chddata c3 on s.code=c3.code ").Where("c1.date", dataResults[60].Date).Where("c2.date", dataResults[30].Date).Where("c3.date", dataResults[0].Date).Where("c1.tclose<c2.tclose and c2.tclose<c3.tclose").Where("s.code in ? AND s.mktcap>20000000000 and s.mktcap<35000000000", codes).Select("s.code,s.name,s.Buy,s.mktcap").Order("c1.turnover desc").Scan(&symbolResults)
		fmt.Println("代码", "名字", "买入价", "总市值")
		for _, v := range symbolResults {
			// 市值格式化为亿

			v.Mktcap = v.Mktcap / 100000000
			fmt.Println(v.Code, v.Name, v.Buy, v.Mktcap)
		}

	}
}
