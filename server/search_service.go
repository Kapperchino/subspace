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
			ID:                space.Space.ID,
			ParentID:          space.Space.ParentID,
			Name:              space.Space.Name,
			Description:       space.Space.Description.String,
			SmallPicture:      getPictureMeta(space.SpaceSmallPicUrl, space.SpaceSmallPicWidth, space.SpaceSmallPicHeight),
			BackgroundPicture: getPictureMeta(space.BackgroundPictureUrl, space.BackgroundPictureWidth, space.BackgroundPictureHeight),
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
			pictures, err := getPicturesForPost(post.Post.ID, queries, c)
			if err != nil {
				return err
			}
			list = append(list, models.Post{
				Id:          post.Post.ID,
				SpaceId:     post.Post.SpaceID.Int64,
				PosterId:    post.Post.PosterID.Int64,
				Topic:       post.Post.Topic.String,
				PosterName:  post.DisplayName,
				ContentType: models.ContentType(post.Post.ContentType.ContentType),
				Body:        post.Post.Body.String, Link: post.Post.Link.String,
				UpVotes:       post.UpVotes,
				DownVotes:     post.DownVotes,
				PosterPicture: getPictureMeta(post.UserPicUrl, post.UserPicWidth, post.UserPicHeight),
				SpacePicture:  getPictureMeta(post.SpaceSmallPicUrl, post.SpaceSmallPicWidth, post.SpaceSmallPicHeight),
				PostPictures:  pictures,
				Created:       post.Post.Created.Time,
				Vote:          getVote(post.IsUpVote, post.VoteType),
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
		pictures, err := getPicturesForPost(post.Post.ID, queries, c)
		if err != nil {
			return err
		}
		list = append(list, models.Post{
			Id:          post.Post.ID,
			SpaceId:     post.Post.SpaceID.Int64,
			PosterId:    post.Post.PosterID.Int64,
			Topic:       post.Post.Topic.String,
			PosterName:  post.DisplayName,
			ContentType: models.ContentType(post.Post.ContentType.ContentType),
			Body:        post.Post.Body.String, Link: post.Post.Link.String,
			UpVotes:       post.UpVotes,
			DownVotes:     post.DownVotes,
			PosterPicture: getPictureMeta(post.UserPicUrl, post.UserPicWidth, post.UserPicHeight),
			SpacePicture:  getPictureMeta(post.SpaceSmallPicUrl, post.SpaceSmallPicWidth, post.SpaceSmallPicHeight),
			PostPictures:  pictures,
			Created:       post.Post.Created.Time,
			Vote:          getVote(post.IsUpVote, post.VoteType),
			SpaceParentId: post.ParentID,
			SpaceName:     post.SpaceName,
		})
	}
	return c.JSON(list)
}
