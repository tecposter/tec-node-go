package mapstructure

import (
	mt "github.com/mitchellh/mapstructure"
)

func Decode(input interface{}, output interface{}) error {
	return mt.Decode(input, output)
}
