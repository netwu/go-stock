package console

import (
	"goravel/app/console/commands"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/schedule"
)

type Kernel struct {
}

func (kernel *Kernel) Schedule() []schedule.Event {
	return []schedule.Event{}
}

func (kernel Kernel) Commands() []console.Command {
	return []console.Command{
		&commands.SendEmails{},
		&commands.Migration{},
		&commands.Start{},
	}
}
