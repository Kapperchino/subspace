package models

type SpaceCreation struct {
	Name        string `validate:"required"`
	Description string `validate:"required"`
	Parent      int64
}

type Space struct {
	ID          int64
	ParentID    int64
	Name        string
	Description string
}

type GetSpace struct {
	Name string
}
