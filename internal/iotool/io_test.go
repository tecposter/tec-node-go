package iotool

import (
	"testing"
	"os"
)

func TestOperateFile(t *testing.T) {
	content := "helloworld, fdfdfief678"
	file := "/tmp/abcdefgfgdffweerfdsa.md"

	WriteFile(file, content)
	got, _ := GetFileContent(file)

	if content != got {
		t.Errorf("Got: %s, Expect: %s", got, content)
	}

	RemoveFile(file)

	if FileExists(file) {
		t.Errorf("file %s should be removed", file)
	}
}

func TestCreateDirIfNotExist(t *testing.T) {
	dir := "/tmp/qwertyuiopasdfghjkl"
	if FileExists(dir) {
		t.Errorf("remove the dir %s first", dir)
	}
	CreateDirIfNotExist(dir)
	if !FileExists(dir) {
		t.Errorf("Create dir %s failed", dir)
	}

	os.Remove(dir)
	if FileExists(dir) {
		t.Errorf("Remove dir %s failed", dir)
	}
}
