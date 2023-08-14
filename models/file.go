package models

type PictureMeta struct {
	Url    string `validate:"required" json:"url"`
	Width  int64  `validate:"required" json:"width"`
	Height int64  `validate:"required" json:"height"`
	Id     int64  `json:"id"`
}

type PictureMetaResult struct {
	Presigned string `json:"presigned"`
	Width     int64  `json:"width"`
	Height    int64  `json:"height"`
	Id        int64  `json:"id"`
	Url       string `json:"url"`
}

type PictureRequestMeta struct {
	Width  int64  `validate:"required" json:"width"`
	Height int64  `validate:"required" json:"height"`
	Url    string `json:"url"`
}

type FileUploadRequest struct {
	PictureMeta *PictureRequestMeta `json:"picture_meta"`
	FileType    FileType            `validate:"required" json:"file_type"`
	IsLink      bool                `json:"is_link"`
}

type FileType string

const (
	FILE_PICTURE FileType = "picture"
	FILE_VIDEO            = "video"
)
