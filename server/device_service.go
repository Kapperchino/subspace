package server

import (
	"database/sql"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/Kapperchino/subspace/util"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type DeviceService struct {
	DB         *sql.DB
	Config     *util.Config
	Validation *util.Validation
}

func (d *DeviceService) getDB() *sql.DB {
	return d.DB
}

func (d *DeviceService) getConfig() *util.Config {
	return d.Config
}

func (d *DeviceService) getValidation() *util.Validation {
	return d.Validation
}

func (d *DeviceService) UpdateRegistration(c *fiber.Ctx) error {
	req := new(models.UpdateDeviceRequest)
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

func (d *DeviceService) GetDevice(c *fiber.Ctx) error {
	userId, _ := c.ParamsInt("userId", -1)
	deviceId := c.Params("deviceId")
	if deviceId == "" || userId == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(d.getDB())
	device, err := queries.GetDeviceByDeviceInfo(c.Context(), gen.GetDeviceByDeviceInfoParams{
		DeviceInfo: sql.NullString{String: deviceId, Valid: true},
		UserID:     sql.NullInt64{Int64: int64(userId), Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.JSON(device)
}
