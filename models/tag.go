package models

type Tag struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type TagName struct {
	Name string `json:"name"`
}
