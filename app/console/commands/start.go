package commands

import (
	"fmt"
	"goravel/app/models"
	"goravel/app/services"
	"sync"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
)

type Start struct {
}

// Signature The name and signature of the console command.
func (receiver *Start) Signature() string {
	return "start"
}

// Description The console command description.
func (receiver *Start) Description() string {
	return `更新股票数据： go run . artisan start
	日期选股，连续三个月上涨，比昨日交易量翻倍， go run . artisan start xg 2021-01-01
	`
}

// Extend The console command extend.
func (receiver *Start) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (receiver *Start) Handle(ctx console.Context) error {
	services.Migrate()
	// action
	name := ctx.Argument(0)
	date := ctx.Argument(1)
	facades.Log.Info(fmt.Sprintf("%s start", name))
	switch name {
	case "symbol":
		services.GetAllStock()
	case "chddata":
		services.GetAllChddataMulity()
	case "bankuai":
		bankuaiService := services.NewBankuaiService()
		bankuaiService.GetAllBankuaiMulity()
	case "xg":
		chddateModel := models.Chddata{}
		chddateModel.Xg(date)
	default:
		getAllData()
	}

	facades.Log.Info(fmt.Sprintf("%s success", name))
	return nil
}

func getAllData() error {
	services.GetAllStock()
	wgAll := sync.WaitGroup{}

	// bankuaiService := services.NewBankuaiService()
	wgAll.Add(1)
	go func() {
		services.GetAllChddataMulity()
		wgAll.Done()
	}()
	// wgAll.Add(1)
	// go func() {
	// bankuaiService.GetAllBankuaiMulity()
	// wgAll.Done()
	// }()
	wgAll.Wait()
	return nil
}
