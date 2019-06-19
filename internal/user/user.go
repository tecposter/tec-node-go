package user

type user struct {
	Uid string `json:"uid"`
	Email string `json:"email"`
	Username string `json:"username"`
	Passhash string `json:"passhash"`
}
