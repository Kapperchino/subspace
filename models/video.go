package models

type VideoProcessingRequest struct {
	Id int64 `json:"id"`
}
type VideoProcessingResponse struct {
	Url       string  `json:"url"`
	Thumbnail string  `json:"thumbnail"`
	Duration  float64 `json:"duration"`
}
