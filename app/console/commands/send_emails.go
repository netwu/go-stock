package commands

import (
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/urfave/cli/v2"
)

type SendEmails struct {
}

//Signature The name and signature of the console command.
func (receiver *SendEmails) Signature() string {
	return "emails"
}

//Description The console command description.
func (receiver *SendEmails) Description() string {
	return "Command description"
}

//Extend The console command extend.
func (receiver *SendEmails) Extend() console.CommandExtend {
	return console.CommandExtend{}
}

//Handle Execute the console command.
func (receiver *SendEmails) Handle(c *cli.Context) error {

	name := c.Args().Get(0)
	email := c.Args().Get(1)
	fmt.Println(name, email)
	return nil
}
