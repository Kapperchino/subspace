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
	if req.FileType == models.FILE_PICTURE {
		return d.uploadPicture(req, queries, c)
	} else if req.FileType == models.FILE_VIDEO {
		return d.uploadVideo(req, queries, c)
	}
	return c.SendStatus(fiber.StatusBadRequest)
}

func (d *FileService) uploadVideo(req *models.FileUploadRequest, queries *gen.Queries, c *fiber.Ctx) error {
	preSigned, fileName, err := d.getPresigned(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting presigned url for picture")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	content := "https://subspaceimg.com/" + fileName
	res, err := queries.CreateVideo(c.Context(), gen.CreateVideoParams{
		Url:          content,
		ProcessState: gen.NullProcessState{ProcessState: gen.ProcessStateOngoing},
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.Status(200).JSON(models.VideoCreationResult{
		Presigned: preSigned,
		Id:        res.ID,
		Url:       content,
	})
}

func (d *FileService) uploadPicture(req *models.FileUploadRequest, queries *gen.Queries, c *fiber.Ctx) error {
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
			Url:       req.PictureMeta.Url,
		})
	}
	preSigned, fileName, err := d.getPresigned(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting presigned url for picture")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if req.PictureMeta != nil {
		content := "https://subspaceimg.com/" + fileName
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
			Url:       content,
		})
	}
	return c.SendStatus(fiber.StatusBadRequest)
}

func (d *FileService) GetFile(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id", -1)
	if id == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(d.getDB())
	pic, err := queries.GetPicture(c.Context(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error().Err(err).Msg("error quering picture")
	}
	res := models.PictureMeta{
		Url:    pic.Url,
		Width:  pic.Width,
		Height: pic.Height,
		Id:     pic.ID,
	}
	return c.JSON(res)
}

func (d *FileService) getPresigned(c *fiber.Ctx) (string, string, error) {
	preSigned, key, err := d.UploadClient.Presign(c.Context())
	return preSigned.URL, key, err
}
