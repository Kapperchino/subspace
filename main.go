package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/go-playground/validator/v10"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Specification struct {
	Debug     bool   `required:"true" default:"false"`
	Port      int    `required:"true" default:"3000"`
	JWTSecret string `required:"true" default:"devSecret"`
}

var validate *validator.Validate
var db *sql.DB
var s Specification

func main() {
	_, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("cant get dir")
	}
	err = envconfig.Process("myapp", &s)
	if err != nil {
		log.Fatal().Err(err)
	}
	app := fiber.New()
	validate = validator.New()
	u, _ := url.Parse("postgres://subspace-dev:devpassword@localhost/joe?sslmode=disable")
	dbm := dbmate.New(u)
	dbm.SchemaFile = "./sql/schema.sql"
	dbm.MigrationsDir = []string{"./sql/migrations"}
	err = dbm.CreateAndMigrate()
	if err != nil {
		log.Fatal().Err(err).Msg("Error with migration")
	}
	dbTemp, err := sql.Open("pgx", "postgres://subspace-dev:devpassword@localhost/joe?sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("Error with connection to db")
	}
	db = dbTemp
	// Create user
	app.Post("/user", createUser)

	// Login route
	app.Post("/login", login)

	// Unauthenticated route
	app.Get("/", accessible)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(s.JWTSecret)},
	}))

	// Restricted Routes
	app.Get("/restricted", restricted)

	app.Listen(":" + strconv.Itoa(s.Port))
}

func validateStruct(input any) error {
	// returns nil or ValidationErrors ( []FieldError )
	err := validate.Struct(input)
	var msgs string
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			msgs += fmt.Sprintf("Field %s needs to be %s", err.Field(), err.Tag())
		}
		// from here you can create your own error messages in whatever language you wish
		return errors.New(msgs)
	}
	return nil
}

func createUser(c *fiber.Ctx) error {
	req := new(models.UserCreation)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err := validateStruct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	salted := hashAndSalt([]byte(req.Password))
	queries := gen.New(db)
	user, err := queries.CreateUser(c.Context(), gen.CreateUserParams{
		Password:    salted,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Bio:         sql.NullString{String: req.Bio},
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	// Create the Claims
	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func login(c *fiber.Ctx) error {
	req := new(models.Login)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	err := validateStruct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	queries := gen.New(db)
	user, err := queries.GetUserFromEmail(c.Context(), req.Email)
	if err != nil {
		log.Error().Err(err).Msg("Error while getting user from db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return c.Status(fiber.StatusBadRequest).SendStatus(400)
	}
	// Create the Claims
	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": t})
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Print(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
