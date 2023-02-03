package services

import "goravel/app/apis/dfcf"

func GetAllChddataMulity() {
	dfcf.NewChddataService().GetAllChddataMulity()
	// dfcf.NewChddataService().Test()
	return
}
