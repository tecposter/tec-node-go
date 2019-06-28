package user

import (
	"github.com/tecposter/tec-node-go/internal/com/dto"
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

func (repo *repository) fetchUser(uid dto.ID) (*user, error) {
	res, err := repo.db.Get(uid)
	if err != nil {
		return nil, err
	}

	var u user
	err = u.unmarshalPair(uid, res)
	return &u, nil
}

func (repo *repository) fetchUIDByEmail(email string) dto.ID {
	res, err := repo.db.Get([]byte(emailPre + email))
	if err != nil {
		return nil
	}

	return res
}

func (repo *repository) hasEmail(email string) bool {
	return repo.fetchUIDByEmail(email) != nil
}

func (repo *repository) fetchUIDByUsername(username string) dto.ID {
	res, err := repo.db.Get([]byte(usernamePre + username))
	if err != nil {
		return nil
	}

	return res
}

func (repo *repository) hasUsername(username string) bool {
	return repo.fetchUIDByUsername(username) != nil
}

var pair = store.Pair
var arr = store.Arr

func (repo *repository) saveUser(uid dto.ID, email string, username string, passhash string) error {
	u := user{
		UID:      uid,
		Email:    email,
		Username: username,
		Passhash: passhash}

	uidKey, data, err := u.marshalPair()
	if err != nil {
		return err
	}

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
