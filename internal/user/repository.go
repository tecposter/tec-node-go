package user

import (
	"encoding/json"
	"github.com/tecposter/tec-node-go/internal/com/store"
)

const (
	emailPre    = "email-"
	usernamePre = "username-"
)

type repository struct {
	db *store.DB
}

func newRepo(userDataDir string) (*repository, error) {
	db, err := store.Open(userDataDir)
	if err != nil {
		return nil, err
	}

	return &repository{db: db}, nil
}

func (repo *repository) Close() {
	repo.db.Close()
}

func (repo *repository) fetchUser(uid string) (*user, error) {
	res, err := repo.db.Get([]byte(uid))
	if err != nil {
		return nil, err
	}

	var u user
	err = json.Unmarshal(res, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (repo *repository) fetchUIDByEmail(email string) string {
	res, err := repo.db.Get([]byte(emailPre + email))
	if err != nil {
		return ""
	}

	return string(res)
}

func (repo *repository) hasEmail(email string) bool {
	return repo.fetchUIDByEmail(email) != ""
}

func (repo *repository) fetchUIDByUsername(username string) string {
	res, err := repo.db.Get([]byte(usernamePre + username))
	if err != nil {
		return ""
	}

	return string(res)
}

func (repo *repository) hasUsername(username string) bool {
	return repo.fetchUIDByUsername(username) != ""
}

var pair = store.Pair
var arr = store.Arr

func (repo *repository) saveUser(uid string, email string, username string, passhash string) error {
	u := user{
		UID:      uid,
		Email:    email,
		Username: username,
		Passhash: passhash}

	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	uidKey := []byte(uid)
	emailKey := []byte(emailPre + email)
	usernameKey := []byte(usernamePre + username)

	err = repo.db.MultiSet(arr(
		pair(uidKey, data),
		pair(emailKey, uidKey),
		pair(usernameKey, uidKey)))
	if err != nil {
		return err
	}

	return nil
}
