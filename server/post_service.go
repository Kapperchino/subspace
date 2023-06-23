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
		Content: sql.NullString{
			String: req.Content,
			Valid:  true,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.JSON(models.Post{
		Id:        post.ID,
		SpaceId:   post.SpaceID.Int64,
		PosterId:  post.PosterID.Int64,
		Topic:     post.Topic,
		Content:   post.Content.String,
		UpVotes:   int64(post.UpVotes.Int32),
		DownVotes: int64(post.DownVotes.Int32),
		Created:   post.Created.Time,
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
		Id:        res.ID,
		SpaceId:   res.SpaceID.Int64,
		PosterId:  res.PosterID.Int64,
		Topic:     res.Topic,
		Content:   res.Content.String,
		UpVotes:   int64(res.UpVotes.Int32),
		DownVotes: int64(res.DownVotes.Int32),
		Created:   res.Created.Time,
	})
}

func (u *PostService) GetSpaces(c *fiber.Ctx) error {
	name := c.Query("name")
	if name == "" {
		parentId := c.QueryInt("parentId")
		// get all spaces
		if parentId == 0 {
			parentId = 1
		}
		queries := gen.New(u.getDB())
		res, err := queries.GetSpacesOfParent(c.Context(), int64(parentId))
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Send(nil)
			}
			log.Error().Err(err).Msg("Error while creating using in db")
			return c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		var spaces []models.Space
		for _, space := range res {
			spaces = append(spaces, models.Space{
				ID:          space.ID,
				ParentID:    space.ParentID,
				Name:        space.Name,
				Description: space.Description.String,
			})
		}
		return c.JSON(spaces)
	}
	queries := gen.New(u.getDB())
	res, err := queries.GetSpaceByName(c.Context(), name)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Send(nil)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.JSON(models.Space{
		ID:          res.ID,
		ParentID:    res.ParentID,
		Name:        res.Name,
		Description: res.Description.String,
	})
}
