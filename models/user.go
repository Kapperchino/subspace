package models

type UserCreation struct {
	Password    string `validate:"required" json:"password"`
	Email       string `validate:"required,email" json:"email"`
	DisplayName string `validate:"required" json:"display_name"`
	Bio         string `json:"bio"`
}

type UserMeta struct {
	UserID      int64  `json:"user_id"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	Token       string `json:"token"`
	Email       string `json:"email"`
}

type UserInfo struct {
	UserID      int64  `json:"user_id"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
}

type Login struct {
	Password string `validate:"required" json:"password"`
	Email    string `validate:"required,email" json:"email"`
}
