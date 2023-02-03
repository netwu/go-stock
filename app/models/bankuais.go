package models

import (
	"fmt"
	"time"

	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
)

// Symbols ...
type Bankuais struct {
	orm.Model
	Id        int       `gorm:"primaryKey"`
	Market    string    `gorm:"type:varchar(255);not null"`
	Code      string    `gorm:"type:varchar(255);not null;index:idx_code,unique"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Bankuai   string    `gorm:"type:varchar(512);not null"`
	Zhuying   string    `gorm:"type:varchar(512);not null"`
	CreatedAt time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP on update current_timestamp;"`
}

func (bankuaiModel *Bankuais) UpdateOrCreate(datas []Bankuais) (err error) {
	for _, data := range datas {
		sql := "INSERT INTO bankuais (symbol,code,name,bankuai,zhuying) VALUES( '%s','%s','%s','%s','%s'  ) ON DUPLICATE KEY UPDATE symbol='%s',code='%s', name='%s',bankuai='%s',zhuying='%s'"
		s := fmt.Sprintf(sql, data.Market, data.Code, data.Name, data.Bankuai, data.Zhuying, data.Market, data.Code, data.Name, data.Bankuai, data.Zhuying)
		facades.Log.Info("bankuaidic sql", s)
		facades.Orm.Query().Exec(s)
		// facades.Orm.Query().Clauses(clause.OnConflict{
		// 	Columns:   []clause.Column{{Name: "code"}},
		// 	DoUpdates: clause.AssignmentColumns([]string{"symbol", "name", "bankuai", "zhuying"}), // column needed to be updated
		// }).Create(&data)

	}

	return
}
