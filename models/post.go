package models

import "time"

type PostCreation struct {
	SpaceId     int64  `validate:"required"`
	PosterId    int64  `validate:"required"`
	Topic       string `validate:"required"`
	Body        string `validate:"required"`
	Content     string
	ContentType ContentType
}

type Post struct {
	Id          int64
	SpaceId     int64
	PosterId    int64
	Topic       string
	Body        string
	Content     string
	ContentType ContentType
	UpVotes     int64
	DownVotes   int64
	Created     time.Time
}

type ContentType string

const (
	CONTENT_VIDEO   ContentType = "video"
	CONTENT_TEXT                = "text"
	CONTENT_PICTURE             = "picture"
)
