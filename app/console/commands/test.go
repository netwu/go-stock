package commands

import (
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/console"
	"github.com/urfave/cli/v2"
)

type Test struct {
}

//Signature The name and signature of the console command.
func (receiver *Test) Signature() string {
	return "test"
}

//Description The console command description.
func (receiver *Test) Description() string {
	return "Test Description"
}

//Extend The console command extend.
func (receiver *Test) Extend() console.CommandExtend {
	return console.CommandExtend{}
}

//Handle Execute the console command.
func (receiver *Test) Handle(c *cli.Context) error {
	for {
		var arr = []int{1, 2, 3, 4, 5}
		for _, item := range arr {
			go receiver.echo(item)
			// break
		}
		time.Sleep(time.Second * 1)
		fmt.Println("\n")

	}

}
func (receiver *Test) echo(i int) {

	time.Sleep(time.Millisecond * time.Duration((6-i)*50))

	fmt.Println(i)
}
