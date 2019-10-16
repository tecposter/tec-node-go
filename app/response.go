package app

// IResponse request interface
type IResponse interface {
	Marshal() ([]byte, error)
}
