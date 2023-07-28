package models

type Device struct {
	Registration string `validate:"required" json:"registration"`
	DeviceID     string `validate:"required" json:"device_id"`
}

type UpdateDeviceRequest struct {
	Registration string `validate:"required" json:"registration"`
	DeviceId     string `validate:"required" json:"device_id"`
	UserId       int64  `validate:"required" json:"userId"`
}
