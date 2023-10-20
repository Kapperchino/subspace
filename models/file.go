package models

type PictureMeta struct {
	Url    string `validate:"required" json:"url"`
	Width  int64  `validate:"required" json:"width"`
	Height int64  `validate:"required" json:"height"`
	Id     int64  `json:"id"`
}

type VideoMeta struct {
	StreamUrl string      `validate:"required" json:"url"`
	Duration  float64     `validate:"required" json:"duration"`
	Id        int64       `validate:"required" json:"id"`
	Thumbnail string      `validate:"required" json:"thumbnail"`
	Status    VideoStatus `validate:"required" json:"status"`
	Width     int64       `validate:"required" json:"width"`
	Height    int64       `validate:"required" json:"height"`
}

type PictureMetaResult struct {
	Presigned string `json:"presigned"`
	Width     int64  `json:"width"`
	Height    int64  `json:"height"`
	Id        int64  `json:"id"`
	Url       string `json:"url"`
}

type VideoCreationResult struct {
	Presigned string `json:"presigned"`
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

type VideoStatus string

const (
	VIDEO_STATUS_FAILED  VideoStatus = "error"
	VIDEO_STATUS_CREATED             = "created"
	VIDEO_STATUS_DONE                = "done"
	VIDEO_STATUS_ONGOING             = "ongoing"
)
