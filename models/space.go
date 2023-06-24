package models

type SpaceCreation struct {
	Name        string `validate:"required" json:"name"`
	Description string `validate:"required" json:"description"`
	Parent      int64  `json:"parent"`
}

type Space struct {
	ID          int64  `json:"id"`
	ParentID    int64  `json:"parent_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
