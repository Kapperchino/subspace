package models

type UserCreation struct {
	Password    string `validate:"required" json:"password"`
	Email       string `validate:"required,email" json:"email"`
	DisplayName string `validate:"required,max=400" json:"display_name"`
	UserAddress string `validate:"required,max=50,alphanum,lowercase" json:"user_address"`
	Bio         string `json:"bio"`
}

type UserBioUpdate struct {
	Bio string `validate:"required" json:"bio"`
}

type UserPictureUpdate struct {
	Id int64 `validate:"required" json:"picture_id"`
}

type UserMeta struct {
	UserID      int64        `json:"user_id"`
	DisplayName string       `json:"display_name"`
	PictureMeta *PictureMeta `json:"picture_meta"`
	Bio         string       `json:"bio"`
	Token       string       `json:"token"`
	Email       string       `json:"email"`
	UserAddress string       `json:"user_address"`
}

type UserInfo struct {
	UserID      int64        `json:"user_id"`
	DisplayName string       `json:"display_name"`
	Bio         string       `json:"bio"`
	PictureMeta *PictureMeta `json:"picture_meta"`
	UserAddress string       `json:"user_address"`
}

type Login struct {
	Password string  `validate:"required" json:"password"`
	Email    string  `validate:"required,email" json:"email"`
	Device   *Device `json:"device"`
}
