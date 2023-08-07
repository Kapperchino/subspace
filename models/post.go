package models

import "time"

type PostCreation struct {
	SpaceId     int64       `validate:"required" json:"space_id"`
	PosterId    int64       `validate:"required" json:"poster_id"`
	Topic       string      `validate:"max=6000" json:"topic"`
	Body        string      `validate:"required_without_all=FileIds Link,max=60000" json:"body"`
	Link        string      `validate:"required_without_all=FileIds Body,max=60000" json:"link"`
	ContentType ContentType `json:"content_type"`
	FileIds     []int64     `json:"file_ids"`
}

type Post struct {
	Id            int64          `json:"id"`
	SpaceId       int64          `json:"space_id"`
	SpacePicture  *PictureMeta   `json:"space_picture"`
	SpaceParentId int64          `json:"space_parent_id"`
	SpaceName     string         `json:"space_name"`
	PosterId      int64          `json:"poster_id"`
	PosterName    string         `json:"poster_name"`
	PosterPicture *PictureMeta   `json:"poster_picture"`
	Topic         string         `json:"topic"`
	Body          string         `json:"body"`
	PostPictures  []*PictureMeta `json:"post_pictures"`
	ContentType   ContentType    `json:"content_type"`
	UpVotes       int64          `json:"up_votes"`
	DownVotes     int64          `json:"down_votes"`
	Created       time.Time      `json:"created"`
	Vote          *Vote          `json:"vote"`
}

type ContentType string

const (
	CONTENT_VIDEO   ContentType = "video"
	CONTENT_TEXT                = "text"
	CONTENT_PICTURE             = "picture"
	CONTENT_LINK                = "link"
)
