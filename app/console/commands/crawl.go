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

type Crawl struct {
}

// Signature The name and signature of the console command.
func (receiver *Crawl) Signature() string {
	return "crawl"
}

// Description The console command description.
func (receiver *Crawl) Description() string {
	return "Command description"
}

// Extend The console command extend.
func (receiver *Crawl) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (receiver *Crawl) Handle(ctx console.Context) error {
	services.Migrate()
	name := ctx.Argument(0)
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
		chddateModel.Xg()
	default:
		getAllData()
	}

	facades.Log.Info(fmt.Sprintf("%s success", name))
	return nil
}

func getAllData() error {
	services.GetAllStock()
	wgAll := sync.WaitGroup{}

	bankuaiService := services.NewBankuaiService()
	wgAll.Add(1)
	go func() {
		services.GetAllChddataMulity()
		wgAll.Done()
	}()
	wgAll.Add(1)
	go func() {
		bankuaiService.GetAllBankuaiMulity()
		wgAll.Done()
	}()
	wgAll.Wait()
	return nil
}
