package draft

import (
	"time"
)

type draft struct {
	PID     string    `json:"pid"`
	Changed time.Time `json:"changed"`
	Cont    content   `json:"cont"`
}

type content struct {
	Typ  string `json:"typ"`
	Body string `json:"body"`
}
