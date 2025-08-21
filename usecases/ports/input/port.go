package input

type Port interface{
	Execute() error
}