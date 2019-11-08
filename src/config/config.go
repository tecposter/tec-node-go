package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"sync"

	"github.com/tecposter/tec-node-go/lib/iotool"
)

type configuration struct {
	File     string `json:"file"`
	BaseDir  string `json:"baseDir"`
	BindAddr string `json:"bindAddr"`

	mux sync.Mutex
}

var c *configuration = &configuration{
	File:     "~/.tec/config.json",
	BaseDir:  "~/.tec",
	BindAddr: "127.0.0.1:7890",
}

var hd string = ""

func homeDir() (string, error) {
	if hd != "" {
		return hd, nil
	}

	hd, err := iotool.HomeDir()
	return hd, err
}

// LoadFromJSONFile load config from json file
func LoadFromJSONFile(file string) error {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mux.Unlock()

	err := parseFilesInConfig()
	if err != nil {
		return err
	}

	parsed, err := parseFile(file)
	if err != nil {
		return err
	}

	if !iotool.FileExists(parsed) {
		// log.Println("config file not exists: ", parsed)

		return errors.New("config file note exists: " + parsed)
	}

	f, err := os.Open(parsed)
	if err != nil {
		return err
	}
	defer f.Close()

	parser := json.NewDecoder(f)
	parser.Decode(c)

	err = parseFilesInConfig()
	if err != nil {
		return err
	}

	return nil
}

func parseFilesInConfig() error {
	var err error
	c.BaseDir, err = parseFile(c.BaseDir)
	if err != nil {
		return err
	}
	c.File, err = parseFile(c.File)
	if err != nil {
		return err
	}

	return nil
}

// BaseDir return base directory
func BaseDir() string {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mux.Unlock()

	return c.BaseDir
}

// BindAddr return application's binding address
func BindAddr() string {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mux.Unlock()

	return c.BindAddr
}

// File returns config file
func File() string {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mux.Unlock()

	return c.File
}

func parseFile(file string) (string, error) {
	if file[0] != '~' {
		return file, nil
	}

	h, err := homeDir()
	if err != nil {
		return file, err
	}
	return path.Join(h, file[1:]), nil
}

/*
myString := "Hello! This is a golangcode.com test ;)"

// Step 1: Convert it to a rune
a := []rune(myString)

// Step 2: Grab the num of chars you need
myShortString := string(a[0:6])

fmt.Println(myShortString)
*/
