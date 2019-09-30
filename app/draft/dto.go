package draft

import (
	"github.com/tecposter/tec-node-go/lib/dto"
	"time"
)

type draft struct {
	ID      dto.ID `json:"pid"`
	Changed int64  `json:"changed"`
	Content string `json:"content"`
}

func newDraft(id dto.ID, content string) *draft {
	return &draft{
		ID:      id,
		Changed: time.Now().UnixNano(),
		Content: content}
}
