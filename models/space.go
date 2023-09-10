package models

type SpaceCreation struct {
	Name                string `validate:"required,max=5000" json:"name"`
	Description         string `validate:"required,max=60000" json:"description"`
	Parent              int64  `json:"parent"`
	SmallPictureId      int64  `json:"small_picture_id"`
	BackgroundPictureId int64  `json:"background_picture_id"`
}

type Space struct {
	ID                int64        `json:"id"`
	ParentID          int64        `json:"parent_id"`
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	SmallPicture      *PictureMeta `json:"small_picture"`
	BackgroundPicture *PictureMeta `json:"background_picture"`
	SubCount          int64        `json:"sub_count"`
}

type SpaceCreationRes struct {
	ID                  int64  `json:"id"`
	ParentID            int64  `json:"parent_id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	SmallPictureId      int64  `json:"small_picture_id"`
	BackgroundPictureId int64  `json:"background_picture_id"`
}
