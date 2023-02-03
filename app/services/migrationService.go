package services

import (
	"goravel/app/models"

	"github.com/goravel/framework/database"
)

func Migrate() {
	gorm, _ := database.NewGormInstance("mysql")
	modelArr := []interface{}{
		&models.Symbols{},
		&models.Bankuais{},
		&models.Chddata{},
		&models.Dfcf{},
		&models.ChddataMonth{},
		&models.BankuaiDic{},
	}
	for _, v := range modelArr {
		if !gorm.Migrator().HasTable(v) {
			gorm.AutoMigrate(v)
		}
	}

}
