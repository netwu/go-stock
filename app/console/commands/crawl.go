package commands

import (
	"fmt"
	"goravel/app/services"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/support/facades"
	"github.com/urfave/cli/v2"
)

type Crawl struct {
}

//Signature The name and signature of the console command.
func (receiver *Crawl) Signature() string {
	return "crawl"
}

//Description The console command description.
func (receiver *Crawl) Description() string {
	return "Crawl Description"
}

//Extend The console command extend.
func (receiver *Crawl) Extend() console.CommandExtend {
	return console.CommandExtend{}
}

//Handle Execute the console command.
func (receiver *Crawl) Handle(c *cli.Context) error {
	name := c.Args().Get(0)
	facades.Log.Info(fmt.Sprintf("%s start", name))
	switch name {
	case "symbol":
		services.GetAllStock()
	case "chddata":
		chddataService := services.NewChddataService()
		chddataService.GetAllChddataMulity()
	default:
		services.GetAllStock()
		chddataService := services.NewChddataService()
		chddataService.GetAllChddataMulity()
	}

	facades.Log.Info(fmt.Sprintf("%s success", name))
	return nil
}
