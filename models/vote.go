package models

type VoteType string

const (
	VOTE_POST    VoteType = "post"
	VOTE_COMMENT          = "comment"
)

type VoteCreation struct {
	UserId          int64    `validate:"required" json:"user_id"`
	PostOrCommentId int64    `validate:"required" json:"post_or_comment_id"`
	IsUpVote        bool     `json:"is_up_vote"`
	VoteType        VoteType `validate:"required" json:"vote_type"`
}

type Vote struct {
	VoteId          int64    `json:"vote_id"`
	UserId          int64    `json:"user_id"`
	PostOrCommentId int64    `json:"post_or_comment_id"`
	IsUpVote        bool     `json:"is_up_vote"`
	VoteType        VoteType `json:"vote_type"`
	IsDeleted       bool     `json:"is_deleted"`
}

type VotesMeta struct {
	VoteId          int64    `json:"vote_id"`
	UpVotes         int64    `json:"up_votes"`
	DownVotes       int64    `json:"down_votes"`
	UserId          int64    `json:"user_id"`
	PostOrCommentId int64    `json:"post_or_comment_id"`
	IsUpVote        bool     `json:"is_up_vote"`
	VoteType        VoteType `json:"vote_type"`
	IsDeleted       bool     `json:"is_deleted"`
}
