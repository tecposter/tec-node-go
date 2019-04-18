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

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
