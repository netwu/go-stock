package commands

import (
	"goravel/app/services"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type SendEmails struct {
}

//Signature The name and signature of the console command.
func (receiver *SendEmails) Signature() string {
	return "migration"
}

//Description The console command description.
func (receiver *SendEmails) Description() string {
	return "Command description"
}

//Extend The console command extend.
func (receiver *SendEmails) Extend() command.Extend {
	return command.Extend{}
}

//Handle Execute the console command.
func (receiver *SendEmails) Handle(ctx console.Context) error {
	services.Migrate()

	return nil
}
