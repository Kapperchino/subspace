package main

import (
	"database/sql"
	"github.com/Kapperchino/subspace/server"
	"github.com/Kapperchino/subspace/util"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/url"
	"strconv"
)

func main() {
	config := util.NewConfig()
	app := fiber.New()
	validation := util.NewValidation()
	u, _ := url.Parse("postgres://subspace-dev:devpassword@localhost/joe?sslmode=disable")
	dbm := dbmate.New(u)
	dbm.SchemaFile = "./sql/schema.sql"
	dbm.MigrationsDir = []string{"./sql/migrations"}
	err := dbm.CreateAndMigrate()
	if err != nil {
		log.Fatal().Err(err).Msg("Error with migration")
	}
	db, err := sql.Open("pgx", "postgres://subspace-dev:devpassword@localhost/joe?sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("Error with connection to db")
	}
	userService := server.UserService{
		DB:         db,
		Config:     config,
		Validation: validation,
	}
	// Create user
	app.Post("/user", userService.CreateUser)

	// Login route
	app.Post("/login", userService.Login)

	// Unauthenticated route
	app.Get("/", accessible)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.JWTSecret)},
	}))

	// Restricted Routes
	app.Get("/restricted", restricted)

	app.Listen(":" + strconv.Itoa(config.Port))
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
