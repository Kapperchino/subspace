package models

type VoteType string

const (
	VOTE_POST    VoteType = "post"
	VOTE_COMMENT          = "comment"
)

type VoteCreation struct {
	UserId          int64 `validate:"required"`
	PostOrCommentId int64 `validate:"required"`
	IsUpVote        bool
	VoteType        VoteType `validate:"required"`
}

type Vote struct {
	VoteId          int64
	UserId          int64
	PostOrCommentId int64
	IsUpVote        bool
	VoteType        VoteType
}
