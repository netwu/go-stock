package models

import (
	"fmt"
	"time"

	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
)

// Symbols ...
type Symbols struct {
	orm.Model
	Id            int       `gorm:"primaryKey"`
	Market        string    `gorm:"type:varchar(255);not null"`
	Code          string    `gorm:"type:varchar(255);not null;index:idx_code,unique"`
	Name          string    `gorm:"type:varchar(255);not null"`
	Trade         float64   `gorm:"type:float;not null"`
	Pricechange   float64   `gorm:"type:float;not null"`
	Changepercent float64   `gorm:"type:float;not null"`
	Buy           float64   `gorm:"type:float;not null"`
	Sell          float64   `gorm:"type:float;not null"`
	Settlement    float64   `gorm:"type:float;comment:昨收;not null"`
	Open          float64   `gorm:"type:float;not null"`
	High          float64   `gorm:"type:float;not null"`
	Low           float64   `gorm:"type:float;not null"`
	Volume        float64   `gorm:"type:float;not null"`
	Amount        float64   `gorm:"type:float;comment:成交额;not null"`
	Ticktime      string    `gorm:"type:varchar(16);not null"`
	Per           float64   `gorm:"type:float;comment:市盈率;not null"`
	Pb            float64   `gorm:"type:float;comment:市净率;not null"`
	Mktcap        float64   `gorm:"type:float;comment:总市值;not null"`
	Nmc           float64   `gorm:"type:float;;comment:流通市值not null"`
	Cbs_rank      float64   `gorm:"type:float;not null"`
	Score         float64   `gorm:"type:float;not null"`
	Content       string    `gorm:"type:LONGTEXT;FULLTEXT INDEX:f_idx_content"`
	Turnoverratio float64   `gorm:"type:float;not null"`
	Roe           float64   `gorm:"type:float;not null"`
	CreatedAt     time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP on update current_timestamp;"`
}

// Amplitude  float64 `gorm:"type:float;comment:振幅;not null"`

func (symbolModel *Symbols) UpdateOrCreate(datas []Symbols) (err error) {

	for _, data := range datas {
		sql := "INSERT INTO symbols (market, code, name, trade, pricechange, changepercent, buy, sell, settlement, open, high, low, volume, amount, ticktime, per, pb, mktcap, nmc, cbs_rank, score, content,turnoverratio, roe) VALUES( '%s', '%s', '%s', %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, '%s', %f, %f, %f, %f, %f, %f, '%s', %f,%f ) ON DUPLICATE KEY UPDATE trade=%f,pricechange=%f,changepercent=%f,buy=%f,sell=%f,settlement=%f,open=%f,high=%f,low=%f,volume=%f,amount=%f,ticktime='%s',per=%f,pb=%f,mktcap=%f,nmc=%f,cbs_rank=%f,score=%f,content='%s',turnoverratio =%f,roe=%f"
		s := fmt.Sprintf(sql, data.Market, data.Code, data.Name, data.Trade, data.Pricechange, data.Changepercent, data.Buy, data.Sell, data.Settlement, data.Open, data.High, data.Low, data.Volume, data.Amount, data.Ticktime, data.Per, data.Pb, data.Mktcap, data.Nmc, data.Cbs_rank, data.Score, data.Content, data.Turnoverratio, data.Roe, data.Trade, data.Pricechange, data.Changepercent, data.Buy, data.Sell, data.Settlement, data.Open, data.High, data.Low, data.Volume, data.Amount, data.Ticktime, data.Per, data.Pb, data.Mktcap, data.Nmc, data.Cbs_rank, data.Score, data.Content, data.Turnoverratio, data.Roe)
		facades.Log.Info("bankuaidic sql", s)
		facades.Orm.Query().Exec(s)

	}

	return
}
