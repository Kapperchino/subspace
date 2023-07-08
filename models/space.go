package models

type SpaceCreation struct {
	Name        string `validate:"required,max=5000" json:"name"`
	Description string `validate:"required,max=60000" json:"description"`
	Parent      int64  `json:"parent"`
	Picture     string `json:"picture"`
}

type Space struct {
	ID          int64  `json:"id"`
	ParentID    int64  `json:"parent_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
}
