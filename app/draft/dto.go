package draft

import (
	"github.com/tecposter/tec-node-go/lib/dto"
	"time"
)

type draftDTO struct {
	ID      dto.ID `json:"id"`
	Changed int64  `json:"changed"`
	Content string `json:"content"`
}

func newDraft(id dto.ID, content string) *draftDTO {
	return &draftDTO{
		ID:      id,
		Changed: time.Now().UnixNano(),
		Content: content}
}
