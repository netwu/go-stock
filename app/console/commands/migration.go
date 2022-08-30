package commands

import (
	"goravel/app/services"

	"github.com/goravel/framework/contracts/console"
	"github.com/urfave/cli/v2"
)

type Migration struct {
}

//Signature The name and signature of the console command.
func (receiver *Migration) Signature() string {
	return "migration"
}

//Description The console command description.
func (receiver *Migration) Description() string {
	return "Migration"
}

//Extend The console command extend.
func (receiver *Migration) Extend() console.CommandExtend {
	return console.CommandExtend{}
}

//Handle Execute the console command.
func (receiver *Migration) Handle(c *cli.Context) error {
	services.Migrate()
	return nil
}
