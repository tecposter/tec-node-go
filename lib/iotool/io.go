package iotool

import (
	"os"
)

// FileExists check file whether exists
func FileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		// log.Fatal(err) bad for transplantation
		panic(err)
	}
}
