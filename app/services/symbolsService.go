package services

import (
	"goravel/app/apis/dfcf"
)

func GetAllStock() error {
	dfcf.GetStock()
	return nil
}
