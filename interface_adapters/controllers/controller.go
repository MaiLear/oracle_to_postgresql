package controllers

import (
	"fmt"
	"gitlab.com/sofia-plus/oracle_to_postgresql/usecases/ports/in"
)

type Controller struct {
	usecase in.Port
}

func (c Controller) Execute() {
	if err := c.usecase.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}
