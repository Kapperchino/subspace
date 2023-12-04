package server

import (
	"database/sql"
	"errors"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/Kapperchino/subspace/util"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type TagService struct {
	DB         *sql.DB
	Config     *util.Config
	Validation *util.Validation
}

func (d *TagService) getDB() *sql.DB {
	return d.DB
}

func (d *TagService) getConfig() *util.Config {
	return d.Config
}

func (d *TagService) getValidation() *util.Validation {
	return d.Validation
}

func (d *TagService) GetPopularTags(c *fiber.Ctx) error {
	queries := gen.New(d.getDB())
	tags, err := queries.GetPopularTags(c.Context(), 7)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var tagRes []models.Tag
	for _, tag := range tags {
		tagRes = append(tagRes, models.Tag{Name: tag.Name, Count: tag.Count})
	}
	return c.JSON(tagRes)
}

func (d *TagService) SearchTags(c *fiber.Ctx) error {
	prefix := c.Query("prefix", "")
	if prefix == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(d.getDB())
	tags, err := queries.SearchPrefixTags(c.Context(), prefix)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var tagRes []models.TagName
	for _, tag := range tags {
		tagRes = append(tagRes, models.TagName{Name: tag.Tag.Name})
	}
	return c.JSON(tagRes)
}
