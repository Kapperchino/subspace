package models

import "time"

type PostCreation struct {
	SpaceId     int64       `validate:"required" json:"space_id"`
	PosterId    int64       `validate:"required" json:"poster_id"`
	Topic       string      `validate:"required,max=6000" json:"topic"`
	Body        string      `validate:"required,max=60000" json:"body"`
	IsUpload    bool        `json:"is_upload"`
	Content     string      `json:"content"`
	ContentType ContentType `json:"content_type"`
}

type Post struct {
	Id            int64       `json:"id"`
	SpaceId       int64       `json:"space_id"`
	SpacePicture  string      `json:"space_picture"`
	SpaceParentId int64       `json:"space_parent_id"`
	SpaceName     string      `json:"space_name"`
	PosterId      int64       `json:"poster_id"`
	PosterName    string      `json:"poster_name"`
	Topic         string      `json:"topic"`
	Body          string      `json:"body"`
	Content       string      `json:"content"`
	IsUpload      bool        `json:"is_upload"`
	Presigned     string      `json:"presigned"`
	ContentType   ContentType `json:"content_type"`
	UpVotes       int64       `json:"up_votes"`
	DownVotes     int64       `json:"down_votes"`
	Created       time.Time   `json:"created"`
	Vote          *Vote       `json:"vote"`
}

type ContentType string

const (
	CONTENT_VIDEO   ContentType = "video"
	CONTENT_TEXT                = "text"
	CONTENT_PICTURE             = "picture"
	CONTENT_LINK                = "link"
)
