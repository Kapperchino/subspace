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

type SubscriptionService struct {
	DB         *sql.DB
	Config     *util.Config
	Validation *util.Validation
}

func (s *SubscriptionService) getDB() *sql.DB {
	return s.DB
}

func (s *SubscriptionService) getConfig() *util.Config {
	return s.Config
}

func (s *SubscriptionService) getValidation() *util.Validation {
	return s.Validation
}

func (s *SubscriptionService) CreateSubscription(c *fiber.Ctx) error {
	req := new(models.SubscriptionCreation)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err := s.getValidation().ValidateStruct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	if req.SpaceId == 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot subscribe to home page")
	}
	tx, err := s.getDB().Begin()
	if err != nil {
		log.Error().Err(err).Msg("Error creating transaction")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer tx.Rollback()
	queries := gen.New(s.getDB()).WithTx(tx)
	existing, err := queries.GetSubscription(c.Context(), gen.GetSubscriptionParams{
		UserID:  req.UserId,
		SpaceID: sql.NullInt64{Int64: req.SpaceId, Valid: true},
	})
	var sub gen.Subscription
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		sub, err = queries.CreateSubscription(c.Context(), gen.CreateSubscriptionParams{
			UserID: req.UserId,
			SpaceID: sql.NullInt64{
				Int64: req.SpaceId,
				Valid: true,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Error while creating using in db")
			return c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(models.Subscription{
			SubscriptionId: sub.ID,
			SpaceId:        sub.SpaceID.Int64,
			UserId:         sub.UserID,
		})
	}
	if err != nil {
		return err
	}
	sub, err = queries.RefreshSubscription(c.Context(), existing.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Error committing transaction")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(models.Subscription{
		SubscriptionId: sub.ID,
		SpaceId:        sub.SpaceID.Int64,
		UserId:         sub.UserID,
	})
}

func (s *SubscriptionService) GetSubscriptionsForUser(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id", -1)
	spaceId := c.QueryInt("spaceId", -1)
	if err != nil && userId == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(s.getDB())
	if spaceId != -1 {
		sub, err := queries.GetActiveSubscription(c.Context(), gen.GetActiveSubscriptionParams{
			UserID:  int64(userId),
			SpaceID: sql.NullInt64{Int64: int64(spaceId), Valid: true},
		})
		if err != nil && errors.Is(err, sql.ErrNoRows) || sub.ID == 0 {
			return c.SendStatus(fiber.StatusNotFound)
		}
		if err != nil {
			log.Error().Err(err).Msg("Error while creating using in db")
			return c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		return c.JSON(models.Subscription{
			UserId:         sub.UserID,
			SpaceId:        sub.SpaceID.Int64,
			SubscriptionId: sub.ID,
		})
	}
	subs, err := queries.GetSubscriptionsForUser(c.Context(), int64(userId))
	if err != nil {
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var res []models.Subscription
	for _, sub := range subs {
		res = append(res, models.Subscription{
			UserId:         sub.UserID,
			SpaceId:        sub.SpaceID.Int64,
			SubscriptionId: sub.ID,
		})
	}
	return c.JSON(res)
}

func (s *SubscriptionService) DeleteSubscription(c *fiber.Ctx) error {
	userId := c.QueryInt("userId", -1)
	spaceId := c.QueryInt("spaceId", -1)
	if userId == -1 || spaceId == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(s.getDB())
	err := queries.DeleteSubscription(c.Context(), gen.DeleteSubscriptionParams{
		SpaceID: sql.NullInt64{Int64: int64(spaceId), Valid: true},
		UserID:  int64(userId),
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.SendStatus(200)
}
