package draft

import (
	"time"
)

type draft struct {
	Pid string `json:"pid"`
	Changed time.Time `json:"changed"`
	Cnt content `json:"cnt"`
}

type content struct {
	Type string `json:"type"`
	Body string `json:"body"`
}
