package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

// Symbols ...
type Dfcf struct {
	orm.Model
	Id        int       `gorm:"primaryKey"`
	Date      string    `gorm:"type:date;not null;uniqueIndex:idx_code_date"`
	Symbol    string    `gorm:"type:varchar(255);not null"`
	Code      string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_code_date"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Zhpj      string    `gorm:"type:varchar(512);comment:综合评价;not null"`
	Scrd      string    `gorm:"type:varchar(512);comment:市场热度;not null"`
	Qsyp      string    `gorm:"type:varchar(512);comment:趋势研判;not null"`
	Jzpg      string    `gorm:"type:varchar(512);comment:价值评估;not null"`
	Zjdx      string    `gorm:"type:varchar(512);comment:资金动向;not null"`
	Score     float64   `gorm:"type:float;comment:评分;not null"`
	Zlkp      string    `gorm:"type:varchar(512);comment:主力控盘;not null"`
	Hypm      int       `gorm:"type:int;comment:行业排名;not null"`
	CreatedAt time.Time `gorm:"type:TIMESTAMP;not null;default:CURRENT_TIMESTAMP"`
}
