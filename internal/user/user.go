package user

import (
	"encoding/json"
	"github.com/tecposter/tec-node-go/internal/com/dto"
)

type user struct {
	UID      dto.ID `json:"uid"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Passhash string `json:"passhash"`
}

func (u *user) marshalPair() ([]byte, []byte, error) {
	uid := u.UID.Bytes()
	arr := [3]string{
		u.Email,
		u.Username,
		u.Passhash}

	data, err := json.Marshal(arr)
	return uid, data, err
}

func (u *user) unmarshalPair(uid, data []byte) error {
	var arr [3]string
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}

	u.UID = dto.ID(uid)
	u.Email = arr[0]
	u.Username = arr[1]
	u.Passhash = arr[2]
	return nil
}
