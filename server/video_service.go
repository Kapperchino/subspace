package server

import (
	"context"
	"database/sql"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/Kapperchino/subspace/util"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"time"
)

type VideoService struct {
	DB               *sql.DB
	Config           *util.Config
	Validation       *util.Validation
	CloudFlareClient *util.CloudFlareClient
}

func (d *VideoService) getDB() *sql.DB {
	return d.DB
}

func (d *VideoService) getConfig() *util.Config {
	return d.Config
}

func (d *VideoService) getValidation() *util.Validation {
	return d.Validation
}

func (d *VideoService) ProcessVideo(c *fiber.Ctx) error {
	req := new(models.VideoProcessingRequest)
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
	video, err := queries.GetVideo(c.Context(), req.Id)
	vid, err := d.CloudFlareClient.Client.StreamUploadFromURL(c.Context(),
		cloudflare.StreamUploadFromURLParameters{
			AccountID:         d.CloudFlareClient.AccountId,
			URL:               video.Url,
			RequireSignedURLs: false,
		})
	if err != nil {
		return err
	}
	err = queries.UpdateVideo(c.Context(), gen.UpdateVideoParams{
		StreamUrl:    sql.NullString{String: vid.Playback.HLS, Valid: true},
		ProcessState: gen.NullProcessState{ProcessState: gen.ProcessStateOngoing, Valid: true},
		Thumbnail:    sql.NullString{String: vid.Thumbnail, Valid: true},
		Duration: sql.NullFloat64{
			Float64: vid.Duration,
			Valid:   true,
		},
		Height: sql.NullInt32{
			Int32: int32(vid.Input.Height),
			Valid: true,
		},
		Width: sql.NullInt32{
			Int32: int32(vid.Input.Width),
			Valid: true,
		},
		ID: req.Id,
	})
	if err != nil {
		return err
	}

	ticker := time.NewTicker(2 * time.Second)
	done := make(chan bool)
	go func() {
		start := time.Now()
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				if t.Sub(start) > 1*time.Minute {
					return
				}
				vid, err := d.CloudFlareClient.Client.StreamGetVideo(context.Background(), cloudflare.StreamParameters{
					AccountID: d.CloudFlareClient.AccountId,
					VideoID:   vid.UID,
				})
				log.Info().Msgf("Polling video api %s", vid.UID)
				if err == nil {
					if vid.ReadyToStream {
						err = queries.UpdateVideo(context.Background(), gen.UpdateVideoParams{
							StreamUrl:    sql.NullString{String: vid.Playback.HLS, Valid: true},
							ProcessState: gen.NullProcessState{ProcessState: gen.ProcessStateDone, Valid: true},
							Thumbnail:    sql.NullString{String: vid.Thumbnail, Valid: true},
							Duration: sql.NullFloat64{
								Float64: vid.Duration,
								Valid:   true,
							},
							Height: sql.NullInt32{
								Int32: int32(vid.Input.Height),
								Valid: true,
							},
							Width: sql.NullInt32{
								Int32: int32(vid.Input.Width),
								Valid: true,
							},
							ID: req.Id,
						})
						log.Info().Msgf("Video %s ready to stream", vid.UID)
						return
					}
				}
			}
		}
	}()

	return c.JSON(models.VideoProcessingResponse{
		Url:       vid.Playback.HLS,
		Thumbnail: vid.Thumbnail,
		Duration:  vid.Duration,
	})
}
