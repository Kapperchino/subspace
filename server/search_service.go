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
			ID:          space.Space.ID,
			ParentID:    space.Space.ParentID,
			Name:        space.Space.Name,
			Description: space.Space.Description.String,
			SmallPicture: &models.PictureMeta{
				Url:    space.Picture.Url,
				Width:  space.Picture.Width,
				Height: space.Picture.Height,
			},
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
			if post.Vote.ID != 0 {
				vote = &models.Vote{
					VoteId:          post.Vote.ID,
					UserId:          post.Vote.UserID,
					PostOrCommentId: post.Vote.PostOrCommentID,
					IsUpVote:        post.Vote.IsUpVote.Bool,
					VoteType:        models.VoteType(post.Vote.VoteType),
					IsDeleted:       post.Vote.IsDeleted.Bool,
				}
			} else {
				vote = nil
			}
			list = append(list, models.Post{
				Id:            post.Post.ID,
				SpaceId:       post.Post.SpaceID.Int64,
				PosterId:      post.Post.PosterID.Int64,
				Topic:         post.Post.Topic.String,
				PosterName:    post.DisplayName,
				ContentType:   models.ContentType(post.Post.ContentType.ContentType),
				Body:          post.Post.Body.String,
				UpVotes:       post.UpVotes,
				DownVotes:     post.DownVotes,
				PosterPicture: getPictureMeta(post.Picture),
				SpacePicture:  getPictureMeta(post.Picture_2),
				PostPicture:   getPictureMeta(post.Picture_3),
				Created:       post.Post.Created.Time,
				Vote:          vote,
				SpaceParentId: post.ParentID,
				SpaceName:     post.SpaceName,
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
		if post.Vote.ID != 0 {
			vote = &models.Vote{
				VoteId:          post.Vote.ID,
				UserId:          post.Vote.UserID,
				PostOrCommentId: post.Vote.PostOrCommentID,
				IsUpVote:        post.Vote.IsUpVote.Bool,
				VoteType:        models.VoteType(post.Vote.VoteType),
				IsDeleted:       post.Vote.IsDeleted.Bool,
			}
		} else {
			vote = nil
		}
		list = append(list, models.Post{
			Id:            post.Post.ID,
			SpaceId:       post.Post.SpaceID.Int64,
			PosterId:      post.Post.PosterID.Int64,
			Topic:         post.Post.Topic.String,
			PosterName:    post.DisplayName,
			ContentType:   models.ContentType(post.Post.ContentType.ContentType),
			Body:          post.Post.Body.String,
			UpVotes:       post.UpVotes,
			DownVotes:     post.DownVotes,
			PosterPicture: getPictureMeta(post.Picture),
			SpacePicture:  getPictureMeta(post.Picture_2),
			PostPicture:   getPictureMeta(post.Picture_3),
			Created:       post.Post.Created.Time,
			Vote:          vote,
			SpaceParentId: post.ParentID,
			SpaceName:     post.SpaceName,
		})
	}
	return c.JSON(list)
}
