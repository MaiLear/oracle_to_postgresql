package controllers

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/sofia-plus/oracle_to_postgresql/usecases/ports/in"
)

type Controller struct {
	usecase in.Port
}

func  NewController(usecase in.Port) Controller{
	return Controller{
		usecase: usecase,
	}
}

func (c Controller) Execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if err := c.usecase.Execute(ctx); err != nil {
		fmt.Println(err.Error())
	}
}
