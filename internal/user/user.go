package user

import (
	"path"
	"sync"
	"errors"
	"os/user"
	"math/rand"
	"fmt"
	"time"
	"context"
	//"github.com/tecposter/tec-server-go/internal/iotool"
)

var u, _ = user.Current()
var baseDir = path.Join(u.HomeDir, ".tec-chain")
var mutex = &sync.Mutex{}

type uidCtnType struct {
	v map[string]string
	mux sync.Mutex
}

var once sync.Once
var uidCtn *uidCtnType

func getUidCtn() *uidCtnType {
	once.Do(func() {
		uidCtn = &uidCtnType{
			v: make(map[string]string)}
	})

	return uidCtn
}

func (ctn *uidCtnType) Get(token string) string {
	if existed, ok := ctn.v[token]; ok {
		return existed
	}

	return ""
}

func (ctn *uidCtnType) Set(token string, uid string) error {
	if _, ok := ctn.v[token]; ok {
		return errors.New(fmt.Sprintf("%s already exists", token))
	}

	ctn.v[token] = uid
	return nil
}

// random string

func initSeed() {
    rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}


func GetUid(token string) string {
	ctn := getUidCtn()
	return ctn.Get(token)
}

func Reg(map[string]interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{});
	return m, nil
}

func Login(_ context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	initSeed()
	token := randStr(15)
	uid := input["uid"].(string)

	getUidCtn().Set(token, uid)

	return map[string]interface{}{"token": token}, nil
}
