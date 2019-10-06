package iotool

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	dirMode = 0755
)

// WriteFile writes content into file
func WriteFile(file string, content string) error {
	tmpfile, err := ioutil.TempFile("", "tec-go")
	if err != nil {
		return err
	}

	//err := ioutil.WriteFile(tmpfile, []byte(content), 0644)
	if _, err := tmpfile.WriteAt([]byte(content), 0); err != nil {
		return err
	}

	err = os.Rename(tmpfile.Name(), file)
	if err != nil {
		return err
	}

	return nil
}

// GetFileContent gets content from file
func GetFileContent(file string) (string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// CurrDir gets directory to current executable
func CurrDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	currDir, err := filepath.Abs(filepath.Dir(ex))
	return currDir, err
}

// MkdirIfNotExist makes directory if the path not exists
func MkdirIfNotExist(path string) error {
	if fileExists(path) {
		return nil
	}

	err := os.MkdirAll(path, dirMode)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		// log.Fatal(err) bad for transplantation
		panic(err)
	}
}
