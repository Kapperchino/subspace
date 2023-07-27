package server

import (
	"database/sql"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/Kapperchino/subspace/util"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"strings"
)

type SearchService struct {
	DB         *sql.DB
	Config     *util.Config
	Validation *util.Validation
}

func (u *SearchService) getDB() *sql.DB {
	return u.DB
}

func (u *SearchService) getConfig() *util.Config {
	return u.Config
}

func (u *SearchService) getValidation() *util.Validation {
	return u.Validation
}

func (u *SearchService) SearchSpace(c *fiber.Ctx) error {
	search := c.Query("term")
	// get all spaces
	if search == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(u.getDB())
	res, err := queries.SearchSpace(c.Context(), search)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Send(nil)
		}
		log.Error().Err(err).Msg("Error while searching db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Space
	for _, space := range res {
		list = append(list, models.Space{
			ID:          space.ID,
			ParentID:    space.ParentID,
			Name:        space.Name,
			Description: space.Description.String,
			Picture:     space.Picture.String,
		})
	}
	return c.JSON(list)
}

func (u *SearchService) SearchPosts(c *fiber.Ctx) error {
	search := c.Query("term")
	userId := c.QueryInt("userId", -1)
	days := c.QueryInt("days", 7)
	isTag := c.QueryBool("isTag", false)
	// get all spaces
	if search == "" || userId == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(u.getDB())
	if isTag {
		res, err := queries.GetPostsWithTagsLatest(c.Context(), gen.GetPostsWithTagsLatestParams{
			UserID: int64(userId),
			Name:   strings.TrimSpace(search),
			Days:   int32(days),
		})
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Send(nil)
			}
			log.Error().Err(err).Msg("Error while searching db")
			return c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		var list []models.Post
		for _, post := range res {
			var vote *models.Vote
			if post.ID_2.Valid {
				vote = &models.Vote{
					VoteId:          post.ID_2.Int64,
					UserId:          post.UserID.Int64,
					PostOrCommentId: post.PostOrCommentID.Int64,
					IsUpVote:        post.IsUpVote.Bool,
					VoteType:        models.VoteType(post.VoteType.VoteType),
					IsDeleted:       post.IsDeleted_2.Bool,
				}
			} else {
				vote = nil
			}
			list = append(list, models.Post{
				Id:            post.ID,
				SpaceId:       post.SpaceID.Int64,
				SpacePicture:  post.SpacePicture.String,
				SpaceParentId: post.ParentID,
				SpaceName:     post.SpaceName,
				PosterId:      post.PosterID.Int64,
				PosterName:    post.DisplayName,
				Topic:         post.Topic.String,
				Body:          post.Body.String,
				Content:       post.Content.String,
				ContentType:   models.ContentType(post.ContentType.ContentType),
				UpVotes:       post.UpVotes,
				DownVotes:     post.DownVotes,
				Created:       post.Created.Time,
				Vote:          vote,
			})
		}
		return c.JSON(list)
	}
	res, err := queries.SearchPost(c.Context(), gen.SearchPostParams{
		UserID:         int64(userId),
		Days:           int32(days),
		PlaintoTsquery: search,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Send(nil)
		}
		log.Error().Err(err).Msg("Error while searching db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Post
	for _, post := range res {
		var vote *models.Vote
		if post.ID_2.Valid {
			vote = &models.Vote{
				VoteId:          post.ID_2.Int64,
				UserId:          post.UserID.Int64,
				PostOrCommentId: post.PostOrCommentID.Int64,
				IsUpVote:        post.IsUpVote.Bool,
				VoteType:        models.VoteType(post.VoteType.VoteType),
				IsDeleted:       post.IsDeleted_2.Bool,
			}
		} else {
			vote = nil
		}
		list = append(list, models.Post{
			Id:            post.ID,
			SpaceId:       post.SpaceID.Int64,
			SpacePicture:  post.SpacePicture.String,
			SpaceParentId: post.ParentID,
			SpaceName:     post.SpaceName,
			PosterId:      post.PosterID.Int64,
			PosterName:    post.DisplayName,
			Topic:         post.Topic.String,
			Body:          post.Body.String,
			Content:       post.Content.String,
			ContentType:   models.ContentType(post.ContentType.ContentType),
			UpVotes:       post.UpVotes,
			DownVotes:     post.DownVotes,
			Created:       post.Created.Time,
			Vote:          vote,
		})
	}
	return c.JSON(list)
}
