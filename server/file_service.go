package server

import (
	"database/sql"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/Kapperchino/subspace/util"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type FileService struct {
	DB         *sql.DB
	Config     *util.Config
	Validation *util.Validation
}

func (d *FileService) getDB() *sql.DB {
	return d.DB
}

func (d *FileService) getConfig() *util.Config {
	return d.Config
}

func (d *FileService) getValidation() *util.Validation {
	return d.Validation
}

func (d *FileService) UploadFile(c *fiber.Ctx) error {
	req := new(models.FileUploadRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err := d.getValidation().ValidateStruct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	queries := gen.New(d.getDB())
	_, err = queries.UpdateDevice(c.Context(), gen.UpdateDeviceParams{
		Registration: sql.NullString{String: req.Registration, Valid: true},
		DeviceInfo:   sql.NullString{String: req.DeviceId, Valid: true},
		UserID:       sql.NullInt64{Int64: req.UserId, Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.SendStatus(fiber.StatusOK)
}
