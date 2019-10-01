package post

import (
	"github.com/tecposter/tec-node-go/lib/dto"
	"time"
)

const (
	typeMarkdown = 1
	typeHTML     = 2
	typeText     = 3
)

type contentDTO struct {
	ID      dto.ID `json:"id"`
	Type    int    `json:"type"`
	Created int64  `json:"created"`
	Content string `json:"content"`
}

func newContent(id dto.ID, contentType int, content string) *contentDTO {
	return &contentDTO{
		ID:      id,
		Type:    contentType,
		Created: time.Now().UnixNano(),
		Content: content}
}
