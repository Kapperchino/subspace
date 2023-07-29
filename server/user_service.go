package server

import (
	"database/sql"
	"errors"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/Kapperchino/subspace/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserService struct {
	DB         *sql.DB
	Config     *util.Config
	Validation *util.Validation
}

func (u *UserService) getDB() *sql.DB {
	return u.DB
}

func (u *UserService) getConfig() *util.Config {
	return u.Config
}

func (u *UserService) getValidation() *util.Validation {
	return u.Validation
}

func (u *UserService) CreateUser(c *fiber.Ctx) error {
	req := new(models.UserCreation)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err := u.getValidation().ValidateStruct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	salted := util.HashAndSalt([]byte(req.Password))
	queries := gen.New(u.getDB())
	user, err := queries.CreateUser(c.Context(), gen.CreateUserParams{
		Password:    salted,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Bio:         sql.NullString{String: req.Bio},
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if strings.Contains(pgErr.Message, "duplicate key value") {
				log.Error().Err(err).Msg("Unique constraint violated")
				if pgErr.ConstraintName == "users_email_key" {
					return c.Status(fiber.StatusConflict).SendString("email")
				} else if pgErr.ConstraintName == "users_display_name_key" {
					return c.Status(fiber.StatusConflict).SendString("display_name")
				}
			}
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// Create the Claims
	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(u.getConfig().JWTSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(models.UserMeta{
		UserID:      user.ID,
		DisplayName: user.DisplayName,
		Bio:         user.Bio.String,
		Token:       t,
		Email:       user.Email,
	})
}

func (u *UserService) Login(c *fiber.Ctx) error {
	req := new(models.Login)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err := u.getValidation().ValidateStruct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	tx, err := u.getDB().Begin()
	if err != nil {
		log.Error().Err(err).Msg("Error creating transaction")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer tx.Rollback()
	queries := gen.New(u.getDB()).WithTx(tx)
	user, err := queries.GetUserFromEmail(c.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error().Err(err).Msg("Error while getting user from db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return c.Status(fiber.StatusBadRequest).SendStatus(400)
	}

	if req.Device != nil {
		device, err := queries.GetDeviceByDeviceInfo(c.Context(), gen.GetDeviceByDeviceInfoParams{
			DeviceInfo: sql.NullString{
				String: req.Device.DeviceID,
				Valid:  true,
			},
			UserID: sql.NullInt64{
				Int64: user.ID,
				Valid: true,
			},
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				_, err := queries.CreateDevice(c.Context(), gen.CreateDeviceParams{
					UserID:       sql.NullInt64{Int64: user.ID, Valid: true},
					Registration: sql.NullString{String: req.Device.Registration, Valid: true},
					DeviceInfo:   sql.NullString{String: req.Device.DeviceID, Valid: true},
				})
				if err != nil {
					log.Error().Err(err).Msg("Error creating device")
					return c.Status(fiber.StatusInternalServerError).SendStatus(500)
				}
			} else {
				log.Error().Err(err).Msg("Error while getting user from db")
				return c.Status(fiber.StatusInternalServerError).SendStatus(500)
			}
		} else {
			_, err = queries.UpdateDevice(c.Context(), gen.UpdateDeviceParams{
				Registration: sql.NullString{String: req.Device.Registration, Valid: true},
				DeviceInfo:   sql.NullString{String: device.DeviceInfo.String, Valid: true},
				UserID:       sql.NullInt64{Int64: user.ID, Valid: true},
			})
			if err != nil {
				log.Error().Err(err).Msg("Error creating device")
				return c.Status(fiber.StatusInternalServerError).SendStatus(500)
			}
		}
	}
	// Create the Claims
	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(u.getConfig().JWTSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	tx.Commit()
	return c.JSON(models.UserMeta{
		UserID:      user.ID,
		DisplayName: user.DisplayName,
		Bio:         user.Bio.String,
		Token:       t,
		Email:       user.Email,
	})
}

func (u *UserService) GetUserById(c *fiber.Ctx) error {
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
	res, err := queries.GetUser(c.Context(), int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Send(nil)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.JSON(models.UserInfo{
		UserID:      res.ID,
		DisplayName: res.DisplayName,
		Bio:         res.Bio.String,
	})
}
