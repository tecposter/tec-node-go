package iotool

import (
	"os/user"
	"testing"
	"fmt"
	"log"
	"path"
)

func TestWriteFile(t *testing.T) {
	WriteFile("/tmp/abcd.md", "helloword")
}

func TestLogin(t *testing.T) {
	//t.Log("test login start")
	usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
	path := path.Join(usr.HomeDir, ".tec-chain")
	fmt.Println(path)
	CreateDirIfNotExist(path)
}
