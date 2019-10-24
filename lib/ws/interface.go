package ws

// IRequest request interface
type IRequest interface {
	CMD() string
	Module() string
	Param(string) (interface{}, bool)
	Marshal() ([]byte, error)
}

// IResponse request interface
type IResponse interface {
	Marshal() ([]byte, error)
	Set(string, interface{})
	SetErr(error)
}
