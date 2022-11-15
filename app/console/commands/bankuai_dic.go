package commands

import (
	"goravel/app/services"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type BankuaiDic struct {
}

//Signature The name and signature of the console command.
func (receiver *BankuaiDic) Signature() string {
	return "bankuaiDic"
}

//Description The console command description.
func (receiver *BankuaiDic) Description() string {
	return "Command description"
}

//Extend The console command extend.
func (receiver *BankuaiDic) Extend() command.Extend {
	return command.Extend{}
}

//Handle Execute the console command.
func (receiver *BankuaiDic) Handle(ctx console.Context) error {
	bankuaidicService := services.NewBankuaiDicService()
	bankuaidicService.PushDicRedis()
	return nil
}
