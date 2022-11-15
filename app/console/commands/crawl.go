package commands

import (
	"fmt"
	"goravel/app/services"
	"sync"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
)

type Crawl struct {
}

//Signature The name and signature of the console command.
func (receiver *Crawl) Signature() string {
	return "crawl"
}

//Description The console command description.
func (receiver *Crawl) Description() string {
	return "Command description"
}

//Extend The console command extend.
func (receiver *Crawl) Extend() command.Extend {
	return command.Extend{}
}

//Handle Execute the console command.
func (receiver *Crawl) Handle(ctx console.Context) error {
	name := ctx.Argument(0)
	facades.Log.Info(fmt.Sprintf("%s start", name))
	switch name {
	case "symbol":
		services.GetAllStock()
	case "chddata":
		chddataService := services.NewChddataService()
		chddataService.GetAllChddataMulity()
	default:
		getAllData()
	}

	facades.Log.Info(fmt.Sprintf("%s success", name))
	return nil
}

func getAllData() error {
	services.GetAllStock()
	wgAll := sync.WaitGroup{}
	wgAll.Add(1)
	go func() {
		chddataService := services.NewChddataService()
		chddataService.GetAllChddataMulity()
		wgAll.Done()
	}()
	wgAll.Add(1)
	go func() {
		bankuaiService := services.NewBankuaiService()
		bankuaiService.GetAllBankuaiMulity()
		wgAll.Done()
	}()
	wgAll.Wait()
	return nil
}
