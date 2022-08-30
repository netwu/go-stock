package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/facades"
	"gorm.io/gorm/clause"
)

// Symbols ...
type Symbols struct {
	orm.Model
	Id            int       `gorm:"primaryKey"`
	Symbol        string    `gorm:"type:varchar(255);not null"`
	Code          string    `gorm:"type:varchar(255);not null;index:idx_code,unique"`
	Name          string    `gorm:"type:varchar(255);not null"`
	Trade         float64   `gorm:"type:float;not null"`
	Pricechange   float64   `gorm:"type:float;not null"`
	Changepercent float64   `gorm:"type:float;not null"`
	Buy           float64   `gorm:"type:float;not null"`
	Sell          float64   `gorm:"type:float;not null"`
	Settlement    float64   `gorm:"type:float;not null"`
	Open          float64   `gorm:"type:float;not null"`
	High          float64   `gorm:"type:float;not null"`
	Low           float64   `gorm:"type:float;not null"`
	Volume        float64   `gorm:"type:float;not null"`
	Amount        float64   `gorm:"type:float;not null"`
	Ticktime      string    `gorm:"type:varchar(16);not null"`
	Per           float64   `gorm:"type:float;not null"`
	Pb            float64   `gorm:"type:float;not null"`
	Mktcap        float64   `gorm:"type:float;not null"`
	Nmc           float64   `gorm:"type:float;not null"`
	Cbs_rank      float64   `gorm:"type:float;not null"`
	Score         float64   `gorm:"type:float;not null"`
	Content       string    `gorm:"type:LONGTEXT;FULLTEXT INDEX:f_idx_content"`
	Turnoverratio float64   `gorm:"type:float;not null"`
	Roe           float64   `gorm:"type:float;not null"`
	CreatedAt     time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP"`
}

func (symbolModel *Symbols) UpdateOrCreate(datas []Symbols) (err error) {
	for _, data := range datas {
		facades.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"trade", "pricechange", "changepercent", "buy", "sell", "settlement", "open", "high", "low", "volume", "amount", "ticktime", "per", "pb", "mktcap", "nmc", "turnoverratio"}), // column needed to be updated
		}).Create(&data)

	}

	return
}
