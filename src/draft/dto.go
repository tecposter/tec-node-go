package draft

import (
	"github.com/tecposter/tec-node-go/lib/dto"
	"time"
)

type draftDTO struct {
	ID      dto.ID `json:"id"`
	Drafted int64  `json:"drafted"`
	Content string `json:"content"`
}

func newDraft(id dto.ID, content string) *draftDTO {
	return &draftDTO{
		ID:      id,
		Drafted: time.Now().UnixNano(),
		Content: content}
}

func (d *draftDTO) Title() string {
	return ""
}

type draftItemDTO struct {
	ID      dto.ID `json:"id"`
	Drafted int64  `json:"drafted"`
	Title   string `json:"title"`
}
