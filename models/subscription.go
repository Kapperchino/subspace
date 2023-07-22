package models

type SubscriptionCreation struct {
	UserId  int64 `validate:"required" json:"user_id"`
	SpaceId int64 `validate:"required" json:"space_id"`
}

type Subscription struct {
	UserId         int64 `validate:"required" json:"user_id"`
	SpaceId        int64 `validate:"required" json:"space_id"`
	SubscriptionId int64 `validate:"required" json:"subscription_id"`
}
