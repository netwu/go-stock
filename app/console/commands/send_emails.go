package commands

import (
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type Migration struct {
}

//Signature The name and signature of the console command.
func (receiver *Migration) Signature() string {
	return "Migration"
}

//Description The console command description.
func (receiver *Migration) Description() string {
	return "Command description"
}

//Extend The console command extend.
func (receiver *Migration) Extend() command.Extend {
	return command.Extend{}
}

//Handle Execute the console command.
func (receiver *Migration) Handle(ctx console.Context) error {
	fmt.Println("hello world")
	return nil
}
