package iotool

import (
	"io/ioutil"
	"os"
	"os/user"
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
	if FileExists(path) {
		return nil
	}

	err := os.MkdirAll(path, dirMode)
	if err != nil {
		return err
	}

	return nil
}

// FileExists check whether file exists
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

// HomeDir returns current user's home direcotry
func HomeDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir, nil
}

/*
func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}
*/
