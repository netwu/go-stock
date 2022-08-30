package services

import (
	"goravel/app/models"

	"github.com/goravel/framework/support/facades"
)

func Migrate() {
	modelArr := []interface{}{
		&models.Symbols{},
		&models.Bankuais{},
		&models.Chddata{},
		&models.Dfcf{},
		&models.Gupings{},
		&models.Guzhis{},
		&models.ChddataMonth{},
	}
	for _, v := range modelArr {
		if !facades.DB.Migrator().HasTable(v) {
			facades.DB.AutoMigrate(v)
		}
	}

}
