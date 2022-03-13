package world

import (
	"fmt"
	"net/rpc"
)


func (c *Character) Add(payload string, reply *string) error {
	test := "gg"
	reply = &test
	return nil
}

func Run() {
	if err := rpc.Register(&Character{}); err != nil {
		fmt.Println(err)
	}
}
