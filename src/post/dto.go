package post

import (
	"github.com/tecposter/tec-node-go/lib/dto"
)

const (
	typeMarkdown = 1
	typeHTML     = 2
	typeText     = 3
)

type commitDTO struct {
	ID        dto.ID `json:"id"`
	PostID    dto.ID `json:"postID"`
	ContentID dto.ID `json:"contentID"`
	Committed int64  `json:"committed"`
}

type postDTO struct {
	ID        dto.ID `json:"id"`
	CommitID  dto.ID `json:"commitID"`
	ContentID dto.ID `json:"contentID"`
	Content   string `json:"content"`
	Posted    int64  `json:"posted"`
	Drafted   int64  `json:"drafted"`
	Committed int64  `json:"committed"`
}

type postItemDTO struct {
	ID        dto.ID `json:"id"`
	CommitID  dto.ID `json:"commitID"`
	ContentID dto.ID `json:"contentID"`
	Title     string `json:"title"`
	Posted    int64  `json:"posted"`
	Drafted   int64  `json:"drafted"`
	Committed int64  `json:"committed"`
}

/*
type contentDTO struct {
	ID      dto.ID `json:"id"`
	Type    int    `json:"type"`
	Content string `json:"content"`
}

func newContent(id dto.ID, contentType int, content string) *contentDTO {
	return &contentDTO{
		ID:      id,
		Type:    contentType,
		Content: content}
}

func newCommit(id dto.ID, postID dto.ID, contentID dto.ID) *commitDTO {
	return &commitDTO{
		ID:        id,
		PostID:    postID,
		ContentID: contentID,
		Created:   time.Now().UnixNano()}
}

type postDTO struct {
	ID       dto.ID `json:"id"`
	CommitID dto.ID `json:"commitID"`
	Created  int64  `json:"created"`
}

func newPost(id dto.ID, commitID dto.ID) *postDTO {
	return &postDTO{
		ID:       id,
		CommitID: commitID,
		Created:  time.Now().UnixNano()}
}
*/
