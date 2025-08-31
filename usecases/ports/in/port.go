package in

import "context"

type Port interface{
	Execute(context.Context) error
}