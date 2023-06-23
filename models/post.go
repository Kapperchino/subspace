package models

import "time"

type PostCreation struct {
	SpaceId  int64  `validate:"required"`
	PosterId int64  `validate:"required"`
	Topic    string `validate:"required"`
	Content  string `validate:"required"`
}

type Post struct {
	Id        int64
	SpaceId   int64
	PosterId  int64
	Topic     string
	Content   string
	UpVotes   int64
	DownVotes int64
	Created   time.Time
}
