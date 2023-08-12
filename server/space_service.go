package server

import (
	"database/sql"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/Kapperchino/subspace/util"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type SpaceService struct {
	DB         *sql.DB
	Config     *util.Config
	Validation *util.Validation
}

func (u *SpaceService) getDB() *sql.DB {
	return u.DB
}

func (u *SpaceService) getConfig() *util.Config {
	return u.Config
}

func (u *SpaceService) getValidation() *util.Validation {
	return u.Validation
}

func (u *SpaceService) CreateSpace(c *fiber.Ctx) error {
	req := new(models.SpaceCreation)
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
	space, err := queries.CreateSpace(c.Context(), gen.CreateSpaceParams{
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		ParentID:    req.Parent,
		SmallPictureID: sql.NullInt64{
			Int64: req.SmallPictureId,
			Valid: req.SmallPictureId != 0,
		},
		BackgroundPictureID: sql.NullInt64{
			Int64: req.BackgroundPictureId,
			Valid: req.BackgroundPictureId != 0,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.JSON(models.SpaceCreationRes{
		ID:                  space.ID,
		ParentID:            space.ParentID,
		Name:                space.Name,
		BackgroundPictureId: space.BackgroundPictureID.Int64,
		SmallPictureId:      space.SmallPictureID.Int64,
		Description:         space.Description.String,
	})
}

func (u *SpaceService) GetSpaceById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id", -1)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	queries := gen.New(u.getDB())
	res, err := queries.GetSpace(c.Context(), int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Send(nil)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.JSON(models.Space{
		ID:                res.Space.ID,
		ParentID:          res.Space.ParentID,
		Name:              res.Space.Name,
		SmallPicture:      getPictureMeta(res.SpaceSmallPicUrl, res.SpaceSmallPicWidth, res.SpaceSmallPicHeight, res.SpaceSmallPictureID.Int64),
		BackgroundPicture: getPictureMeta(res.BackgroundPictureUrl, res.BackgroundPictureWidth, res.BackgroundPictureHeight, res.BackgroundPictureID.Int64),
		Description:       res.Space.Description.String,
	})
}

func (u *SpaceService) GetSpaces(c *fiber.Ctx) error {
	name := c.Query("name")
	parentId := c.QueryInt("parentId")
	// get all spaces
	if parentId == 0 {
		parentId = 1
	}
	if name == "" {
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
				ID:                space.Space.ID,
				ParentID:          space.Space.ParentID,
				Name:              space.Space.Name,
				SmallPicture:      getPictureMeta(space.SpaceSmallPicUrl, space.SpaceSmallPicWidth, space.SpaceSmallPicHeight, space.SpaceSmallPictureID.Int64),
				BackgroundPicture: getPictureMeta(space.BackgroundPictureUrl, space.BackgroundPictureWidth, space.BackgroundPictureHeight, space.BackgroundPictureID.Int64),
				Description:       space.Space.Description.String,
			})
		}
		return c.JSON(spaces)
	}
	queries := gen.New(u.getDB())
	res, err := queries.GetSpaceByName(c.Context(), gen.GetSpaceByNameParams{
		Name:     name,
		ParentID: int64(parentId),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Send(nil)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.JSON(models.Space{
		ID:                res.Space.ID,
		ParentID:          res.Space.ParentID,
		Name:              res.Space.Name,
		SmallPicture:      getPictureMeta(res.SpaceSmallPicUrl, res.SpaceSmallPicWidth, res.SpaceSmallPicHeight, res.SpaceSmallPictureID.Int64),
		BackgroundPicture: getPictureMeta(res.BackgroundPictureUrl, res.BackgroundPictureWidth, res.BackgroundPictureHeight, res.BackgroundPictureID.Int64),
		Description:       res.Space.Description.String,
	})
}

func (u *SpaceService) GetSpacesForUser(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id", -1)
	// get all spaces
	if userId == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(u.getDB())
	res, err := queries.GetUserSpaces(c.Context(), int64(userId))
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
			ID:                space.Space.ID,
			ParentID:          space.Space.ParentID,
			Name:              space.Space.Name,
			SmallPicture:      getPictureMeta(space.SpaceSmallPicUrl, space.SpaceSmallPicWidth, space.SpaceSmallPicHeight, space.SpaceSmallPictureID.Int64),
			BackgroundPicture: getPictureMeta(space.BackgroundPictureUrl, space.BackgroundPictureWidth, space.BackgroundPictureHeight, space.BackgroundPictureID.Int64),
			Description:       space.Space.Description.String,
		})
	}
	return c.JSON(spaces)
}
