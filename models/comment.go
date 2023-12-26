package models

import "time"

type CommentCreation struct {
	PosterId    int64       `validate:"required" json:"poster_id"`
	ParentId    int64       `json:"parent_id"`
	PostId      int64       `validate:"required" json:"post_id"`
	Body        string      `validate:"required,max=40000" json:"body"`
	ContentType ContentType `json:"content_type"`
	FileIds     []int64     `json:"file_ids"`
}

type Comment struct {
	Id              int64          `json:"id"`
	PosterId        int64          `json:"poster_id"`
	PostId          int64          `json:"post_id"`
	ParentId        int64          `json:"parent_id"`
	PosterName      string         `json:"poster_name"`
	Body            string         `json:"body"`
	ContentType     ContentType    `json:"content_type"`
	CommentPictures []*PictureMeta `json:"comment_pictures"`
	CommentVideos   []*VideoMeta   `json:"comment_videos"`
	UpVotes         int64          `json:"up_votes"`
	DownVotes       int64          `json:"down_votes"`
	Created         time.Time      `json:"created"`
	Vote            *Vote          `json:"vote"`
	PosterPicture   *PictureMeta   `json:"poster_picture"`
}
