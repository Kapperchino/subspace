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
	DB           *sql.DB
	Config       *util.Config
	Validation   *util.Validation
	UploadClient *util.UploadClient
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
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(d.getDB())
	//only supporting file atm
	if req.FileType != models.FILE_PICTURE {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if req.IsLink {
		res, err := queries.CreatePicture(c.Context(), gen.CreatePictureParams{
			Url:    req.PictureMeta.Url,
			Width:  req.PictureMeta.Width,
			Height: req.PictureMeta.Height,
		})
		if err != nil {
			log.Error().Err(err).Msg("Error while creating using in db")
			return c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		return c.Status(200).JSON(models.PictureMetaResult{
			Presigned: "",
			Width:     res.Width,
			Height:    res.Height,
			Id:        res.ID,
		})
	}
	preSigned, fileName, err := d.getPresigned(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting presigned url for picture")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if req.PictureMeta != nil {
		content := "https://pub-cab547f3a0034c6083d1d10ab8298a3f.r2.dev/" + fileName
		res, err := queries.CreatePicture(c.Context(), gen.CreatePictureParams{
			Url:    content,
			Width:  req.PictureMeta.Width,
			Height: req.PictureMeta.Height,
		})
		if err != nil {
			log.Error().Err(err).Msg("Error while creating using in db")
			return c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		return c.Status(200).JSON(models.PictureMetaResult{
			Presigned: preSigned,
			Width:     res.Width,
			Height:    res.Height,
			Id:        res.ID,
		})
	}
	return c.SendStatus(fiber.StatusBadRequest)
}

func (d *FileService) getPresigned(c *fiber.Ctx) (string, string, error) {
	preSigned, key, err := d.UploadClient.Presign(c.Context())
	return preSigned.URL, key, err
}
