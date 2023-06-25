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
	voteExists := err == nil
	//row exists
	if voteExists {
		if vote.IsDeleted.Bool {
			err = queries.RefreshVote(c.Context(), vote.ID)
			if err != nil {
				log.Error().Err(err).Msg("Error creating transaction")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		} else if req.IsUpVote == vote.IsUpVote.Bool {
			err = queries.DeleteVote(c.Context(), vote.ID)
			if err != nil {
				log.Error().Err(err).Msg("Error creating transaction")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		} else {
			err = queries.SetVote(c.Context(), gen.SetVoteParams{
				IsUpVote: sql.NullBool{
					Bool:  req.IsUpVote,
					Valid: true,
				},
				ID: vote.ID,
			})
			if err != nil {
				log.Error().Err(err).Msg("Error creating transaction")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}
	} else {
		vote, err = queries.CreateVote(c.Context(), gen.CreateVoteParams{
			UserID:          req.UserId,
			PostOrCommentID: req.PostOrCommentId,
			VoteType:        gen.VoteType(req.VoteType),
			IsUpVote: sql.NullBool{
				Bool:  req.IsUpVote,
				Valid: true},
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
