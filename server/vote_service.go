package server

import (
	"database/sql"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/Kapperchino/subspace/util"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type VoteService struct {
	DB         *sql.DB
	Config     *util.Config
	Validation *util.Validation
}

func (u *VoteService) getDB() *sql.DB {
	return u.DB
}

func (u *VoteService) getConfig() *util.Config {
	return u.Config
}

func (u *VoteService) getValidation() *util.Validation {
	return u.Validation
}

// CreateVote votes/
func (u *VoteService) CreateVote(c *fiber.Ctx) error {
	req := new(models.VoteCreation)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err := u.getValidation().ValidateStruct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	tx, err := u.getDB().Begin()
	if err != nil {
		log.Error().Err(err).Msg("Error creating transaction")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer tx.Rollback()
	queries := gen.New(u.getDB()).WithTx(tx)
	vote, err := queries.GetVoteForPostOrCommentForUser(c.Context(),
		gen.GetVoteForPostOrCommentForUserParams{
			UserID:          req.UserId,
			PostOrCommentID: req.PostOrCommentId,
		})
	if err != nil && err != sql.ErrNoRows {
		log.Error().Err(err).Msg("Error creating transaction")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	rowExists := err == nil
	//row exists
	if rowExists {
		//changing vote
		if vote.IsUpVote.Bool != req.IsUpVote {
			err := queries.SetVote(c.Context(), gen.SetVoteParams{
				IsUpVote: sql.NullBool{Bool: req.IsUpVote, Valid: true},
				ID:       vote.ID,
			})
			if err != nil {
				log.Error().Err(err).Msg("Error creating transaction")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			if req.IsUpVote {
				//changing from downVote to upVote
				err = u.deDownVoteCommentOrPost(c, req, queries)
				if err != nil {
					log.Error().Err(err).Msg("Error creating transaction")
					return c.SendStatus(fiber.StatusInternalServerError)
				}
				//upvote the counter
				err := u.upvoteCommentOrPost(c, req, queries)
				if err != nil {
					log.Error().Err(err).Msg("Error creating transaction")
					return c.SendStatus(fiber.StatusInternalServerError)
				}
			} else {
				//changing from upvote to downvote
				err := u.deUpvoteCommentOrPost(c, req, queries)
				if err != nil {
					log.Error().Err(err).Msg("Error creating transaction")
					return c.SendStatus(fiber.StatusInternalServerError)
				}
				//downVote the counter
				err = u.downVoteCommentOrPost(c, req, queries)
				if err != nil {
					log.Error().Err(err).Msg("Error creating transaction")
					return c.SendStatus(fiber.StatusInternalServerError)
				}
			}
		} else if !vote.IsDeleted.Bool {
			//unliking or undisliking a vote
			if vote.IsUpVote.Bool {
				err := u.deUpvoteCommentOrPost(c, req, queries)
				if err != nil {
					log.Error().Err(err).Msg("Error creating transaction")
					return c.SendStatus(fiber.StatusInternalServerError)
				}
			} else {
				err := u.deDownVoteCommentOrPost(c, req, queries)
				if err != nil {
					log.Error().Err(err).Msg("Error creating transaction")
					return c.SendStatus(fiber.StatusInternalServerError)
				}
			}
			err := queries.DeleteVote(c.Context(), vote.ID)
			if err != nil {
				log.Error().Err(err).Msg("Error creating transaction")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		} else {
			err := queries.RefreshVote(c.Context(), vote.ID)
			if err != nil {
				log.Error().Err(err).Msg("Error creating transaction")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}
	} else {
		if req.IsUpVote {
			err := u.upvoteCommentOrPost(c, req, queries)
			if err != nil {
				log.Error().Err(err).Msg("Error creating transaction")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		} else {
			err := u.downVoteCommentOrPost(c, req, queries)
			if err != nil {
				log.Error().Err(err).Msg("Error creating transaction")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}
	}

	if !rowExists {
		vote, err = queries.CreateVote(c.Context(), gen.CreateVoteParams{
			UserID:          req.UserId,
			PostOrCommentID: req.PostOrCommentId,
			IsUpVote: sql.NullBool{
				Bool:  req.IsUpVote,
				Valid: true,
			},
			VoteType: gen.VoteType(req.VoteType),
		})
		if err != nil {
			log.Error().Err(err).Msg("Error creating transaction")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Error committing transaction")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(models.Vote{
		VoteId:          vote.ID,
		UserId:          vote.UserID,
		PostOrCommentId: vote.PostOrCommentID,
		VoteType:        models.VoteType(vote.VoteType),
		IsUpVote:        req.IsUpVote,
	})
}

func (u *VoteService) upvoteCommentOrPost(c *fiber.Ctx, req *models.VoteCreation, queries *gen.Queries) error {
	var err error
	if req.VoteType == models.VOTE_POST {
		err = queries.UpVotePost(c.Context(), req.PostOrCommentId)
	} else {
		err = queries.UpVoteComment(c.Context(), req.PostOrCommentId)
	}
	if err != nil {
		log.Error().Err(err).Msg("Error creating transaction")
		return err
	}
	return nil
}

func (u *VoteService) deUpvoteCommentOrPost(c *fiber.Ctx, req *models.VoteCreation, queries *gen.Queries) error {
	var err error
	if req.VoteType == models.VOTE_POST {
		err = queries.DeUpVotePost(c.Context(), req.PostOrCommentId)
	} else {
		err = queries.DeUpVoteComment(c.Context(), req.PostOrCommentId)
	}
	if err != nil {
		log.Error().Err(err).Msg("Error creating transaction")
		return err
	}
	return nil
}

func (u *VoteService) deDownVoteCommentOrPost(c *fiber.Ctx, req *models.VoteCreation, queries *gen.Queries) error {
	var err error
	if req.VoteType == models.VOTE_POST {
		err = queries.DeDownVotePost(c.Context(), req.PostOrCommentId)
	} else {
		err = queries.DeDownVoteComment(c.Context(), req.PostOrCommentId)
	}
	if err != nil {
		log.Error().Err(err).Msg("Error creating transaction")
		return err
	}
	return nil
}

func (u *VoteService) downVoteCommentOrPost(c *fiber.Ctx, req *models.VoteCreation, queries *gen.Queries) error {
	var err error
	if req.VoteType == models.VOTE_POST {
		err = queries.DownVotePost(c.Context(), req.PostOrCommentId)
	} else {
		err = queries.DownVoteComment(c.Context(), req.PostOrCommentId)
	}
	if err != nil {
		log.Error().Err(err).Msg("Error creating transaction")
		return err
	}
	return nil
}
