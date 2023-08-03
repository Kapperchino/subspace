package models

type PictureMeta struct {
	Url    string ` json:"url"`
	Width  int64  `validate:"required" json:"width"`
	Height int64  `validate:"required" json:"height"`
}

type PictureMetaResult struct {
	Url    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
	Id     int64  `json:"id"`
}

type FileUploadRequest struct {
	PictureMeta *PictureMeta `json:"picture_meta"`
	FileType    FileType     `validate:"required" json:"file_type"`
}

type FileType string

const (
	FILE_PICTURE FileType = "picture"
	FILE_VIDEO            = "video"
)
