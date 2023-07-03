package server

import (
	"database/sql"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/Kapperchino/subspace/util"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type PostService struct {
	DB         *sql.DB
	Config     *util.Config
	Validation *util.Validation
}

func (u *PostService) getDB() *sql.DB {
	return u.DB
}

func (u *PostService) getConfig() *util.Config {
	return u.Config
}

func (u *PostService) getValidation() *util.Validation {
	return u.Validation
}

func (u *PostService) CreatePost(c *fiber.Ctx) error {
	req := new(models.PostCreation)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err := u.getValidation().ValidateStruct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	queries := gen.New(u.getDB())
	if req.ContentType == "" {
		req.ContentType = models.CONTENT_TEXT
	}
	post, err := queries.CreatePost(c.Context(), gen.CreatePostParams{
		SpaceID: sql.NullInt64{
			Int64: req.SpaceId,
			Valid: true,
		},
		PosterID: sql.NullInt64{
			Int64: req.PosterId,
			Valid: true,
		},
		Topic: req.Topic,
		Body: sql.NullString{
			String: req.Body,
			Valid:  true,
		},
		Content: sql.NullString{
			String: req.Content,
			Valid:  req.Content != "",
		},
		ContentType: gen.NullContentType{
			ContentType: gen.ContentType(req.ContentType),
			Valid:       req.ContentType != "",
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.JSON(models.Post{
		Id:          post.ID,
		SpaceId:     post.SpaceID.Int64,
		PosterId:    post.PosterID.Int64,
		Body:        post.Body.String,
		Topic:       post.Topic,
		Content:     post.Content.String,
		ContentType: models.ContentType(post.ContentType.ContentType),
		UpVotes:     0,
		DownVotes:   0,
		Created:     post.Created.Time,
	})
}

func (u *PostService) GetPostById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id", -1)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if id == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(u.getDB())
	res, err := queries.GetPost(c.Context(), int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Send(nil)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.JSON(models.Post{
		Id:           res.ID,
		SpaceId:      res.SpaceID.Int64,
		PosterId:     res.PosterID.Int64,
		PosterName:   res.DisplayName,
		SpacePicture: res.SpacePicture.String,
		Topic:        res.Topic,
		Body:         res.Body.String,
		ContentType:  models.ContentType(res.ContentType.ContentType),
		Content:      res.Content.String,
		UpVotes:      res.UpVotes,
		DownVotes:    res.DownVotes,
		Created:      res.Created.Time,
	})
}

func (u *PostService) GetPosts(c *fiber.Ctx) error {
	spaceId := c.QueryInt("spaceId", -1)
	userId := c.QueryInt("userId", -1)
	if spaceId != -1 && userId == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	//TODO:get post by popularity
	if spaceId == -1 {
		spaceId = 1
	}
	queries := gen.New(u.getDB())
	if userId != -1 {
		res, err := queries.GetPostsForUser(c.Context(), sql.NullInt64{Int64: int64(userId), Valid: true})
		if err != nil {
			if err == sql.ErrNoRows {
				return c.SendStatus(fiber.StatusOK)
			}
			log.Error().Err(err).Msg("Error while creating using in db")
			return c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		var list []models.Post
		for _, post := range res {
			list = append(list, models.Post{
				Id:           post.ID,
				SpaceId:      post.SpaceID.Int64,
				PosterId:     post.PosterID.Int64,
				SpacePicture: post.SpacePicture.String,
				Topic:        post.Topic,
				Content:      post.Content.String,
				PosterName:   post.DisplayName,
				ContentType:  models.ContentType(post.ContentType.ContentType),
				Body:         post.Body.String,
				UpVotes:      post.UpVotes,
				DownVotes:    post.DownVotes,
				Created:      post.Created.Time,
			})
		}
		return c.JSON(list)
	}
	res, err := queries.GetPostsForSpace(c.Context(), sql.NullInt64{Int64: int64(spaceId), Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.SendStatus(fiber.StatusOK)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Post
	for _, post := range res {
		list = append(list, models.Post{
			Id:           post.ID,
			SpaceId:      post.SpaceID.Int64,
			PosterId:     post.PosterID.Int64,
			SpacePicture: post.SpacePicture.String,
			Topic:        post.Topic,
			Content:      post.Content.String,
			PosterName:   post.DisplayName,
			ContentType:  models.ContentType(post.ContentType.ContentType),
			Body:         post.Body.String,
			UpVotes:      post.UpVotes,
			DownVotes:    post.DownVotes,
			Created:      post.Created.Time,
		})
	}
	return c.JSON(list)
}
