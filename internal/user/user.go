package user

type user struct {
	UID      string `json:"uid"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Passhash string `json:"passhash"`
}
