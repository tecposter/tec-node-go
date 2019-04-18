package ipfs

import (
	"encoding/json"
	"bytes"
	"strings"
	"sync"
	shell "github.com/ipfs/go-ipfs-api"
)

// https://godoc.org/github.com/ipfs/go-ipfs-api

var sh = shell.NewShell("localhost:5001")
var mutex = &sync.Mutex{}
var latestCid string

func Add(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	mutex.Lock()
	cid, err := sh.Add(bytes.NewReader(b))
	latestCid = cid
	mutex.Unlock()

	return cid, err
}

func AddStr(s string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	return sh.Add(strings.NewReader(s))
}



func RecvByCid(cid string, out interface{}) error {
	//func (s *Shell) Cat(path string) (io.ReadCloser, error)
	mutex.Lock()
	ioRes, err := sh.Cat("/ipfs/" + cid)
	mutex.Unlock()

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
    buf.ReadFrom(ioRes)

	return json.Unmarshal(buf.Bytes(), out)
}

func FetchStr(cid string) (string, error) {
	mutex.Lock()
	io, err := sh.Cat("/ipfs/" + cid)
	mutex.Unlock()

	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(io)

	return buf.String(), nil
}

func LatestCid() string {
	return latestCid
}
