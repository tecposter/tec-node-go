package iotool

import (
	"os"
	"io/ioutil"
)

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

func GetFileContent(file string) (string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func RemoveFile(file string) error {
	return os.Remove(file)
}

func CreateDirIfNotExist(dir string) error {
	if FileExists(dir) {
		return nil
	}

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	return nil
	/*
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
	*/
}

func FileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}

	return true
}
