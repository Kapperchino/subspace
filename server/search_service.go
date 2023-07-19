package server

import (
	"database/sql"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/Kapperchino/subspace/util"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
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
	search := c.Params("spaceName")
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
