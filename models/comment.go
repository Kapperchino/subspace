package models

import "time"

type CommentCreation struct {
	PosterId    int64       `validate:"required" json:"poster_id"`
	ParentId    int64       `json:"parent_id"`
	Body        string      `validate:"required" json:"body"`
	Content     string      `json:"content"`
	ContentType ContentType `json:"content_type"`
}

type Comment struct {
	Id          int64       `json:"id"`
	PosterId    int64       `json:"poster_id"`
	ParentId    int64       `json:"parent_id"`
	PosterName  string      `json:"poster_name"`
	Body        string      `json:"body"`
	Content     string      `json:"content"`
	ContentType ContentType `json:"content_type"`
	UpVotes     int64       `json:"up_votes"`
	DownVotes   int64       `json:"down_votes"`
	Created     time.Time   `json:"created"`
}
