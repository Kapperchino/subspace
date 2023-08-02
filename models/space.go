package models

type SpaceCreation struct {
	Name              string       `validate:"required,max=5000" json:"name"`
	Description       string       `validate:"required,max=60000" json:"description"`
	Parent            int64        `json:"parent"`
	SmallPicture      *PictureMeta `json:"small_picture"`
	BackgroundPicture *PictureMeta `json:"background_picture"`
}

type Space struct {
	ID           int64        `json:"id"`
	ParentID     int64        `json:"parent_id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	SmallPicture *PictureMeta `json:"small_picture"`
}
