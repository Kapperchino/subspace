package main

import (
	"database/sql"
	"github.com/Kapperchino/subspace/server"
	"github.com/Kapperchino/subspace/util"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/goccy/go-json"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	config := util.NewConfig()
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	validation := util.NewValidation()
	validation.Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
	u, _ := url.Parse(config.DatabaseUrl)
	dbm := dbmate.New(u)
	dbm.SchemaFile = "./sql/schema.sql"
	dbm.MigrationsDir = []string{"./sql/migrations"}
	err := dbm.CreateAndMigrate()
	if err != nil {
		log.Fatal().Err(err).Msg("Error with migration")
	}
	db, err := sql.Open("pgx", config.DatabaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("Error with connection to db")
	}
	uploadClient, err := util.NewUploadClient(config.BucketName, config.BucketAccountId, config.BucketKeyId, config.BucketKeySecret)
	if err != nil {
		log.Fatal().Err(err).Msg("Error with connection to objectStore")
	}
	userService := server.UserService{
		DB:         db,
		Config:     config,
		Validation: validation,
	}

	spaceService := server.SpaceService{
		DB:         db,
		Config:     config,
		Validation: validation,
	}

	postService := server.PostService{
		DB:         db,
		Config:     config,
		Validation: validation,
	}

	voteService := server.VoteService{
		DB:         db,
		Config:     config,
		Validation: validation,
	}

	commentService := server.CommentService{
		DB:         db,
		Config:     config,
		Validation: validation,
	}

	subscriptionService := server.SubscriptionService{
		DB:         db,
		Config:     config,
		Validation: validation,
	}

	searchService := server.SearchService{
		DB:         db,
		Config:     config,
		Validation: validation,
	}

	deviceService := server.DeviceService{
		DB:         db,
		Config:     config,
		Validation: validation,
	}

	fileService := server.FileService{
		DB:           db,
		Config:       config,
		Validation:   validation,
		UploadClient: uploadClient,
	}

	tagService := server.TagService{
		DB:         db,
		Config:     config,
		Validation: validation,
	}

	app.Use(cors.New())
	// Create user
	app.Post("/auth/user", userService.CreateUser)

	// Login route
	app.Post("/auth/login", userService.Login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.JWTSecret)},
	}))

	// Restricted Routes
	//user
	app.Get("/users/:id", userService.GetUserById)

	app.Put("/users/:id", userService.UpdateUserBio)

	app.Put("/users/:id/picture", userService.UpdatePicture)

	//spaces
	app.Post("/spaces", spaceService.CreateSpace)

	app.Get("/spaces/", spaceService.GetSpaces)

	app.Get("/spaces/sort/", spaceService.GetSpacesSorted)

	app.Get("/spaces/:id", spaceService.GetSpaceById)

	app.Get("/spaces/users/:id", spaceService.GetSpacesForUser)

	//posts
	app.Post("/posts", postService.CreatePost)

	app.Get("/posts/:id", postService.GetPostById)

	app.Get("/posts/spaces/:id", postService.GetPosts)

	app.Get("/posts/users/:id/subscriptions", postService.GetPostsForSubscription)

	app.Get("/posts/", postService.GetPostsByName)

	app.Get("/posts/tags/:name", postService.GetPostsForTag)

	app.Get("/posts/users/:id", postService.GetPostsForUser)

	//votes
	app.Post("/votes", voteService.CreateVote)

	app.Get("/votes/:id", voteService.GetVote)

	app.Get("/votes/posts/:id", voteService.GetVotesForPost)

	app.Get("/votes/comments/:id", voteService.GetVotesForComment)

	//comments
	app.Post("/comments", commentService.CreateComment)

	app.Get("/comments/:id", commentService.GetCommentById)

	app.Get("/comments/", commentService.GetComments)

	//subscription
	app.Put("/subscriptions/", subscriptionService.CreateSubscription)

	app.Get("/subscriptions/users/:id", subscriptionService.GetSubscriptionsForUser)

	app.Delete("/subscriptions/", subscriptionService.DeleteSubscription)

	//search
	app.Get("/search/spaces/", searchService.SearchSpace)

	app.Get("/search/posts/", searchService.SearchPosts)

	app.Get("/search/users/", userService.SearchUsers)

	//device
	app.Put("/devices", deviceService.UpdateRegistration)

	app.Get("/devices/users/:", deviceService.GetDevice)

	//files
	app.Put("/files", fileService.UploadFile)

	app.Get("/files/:id", fileService.GetFile)

	//tags
	app.Get("/tags", tagService.GetPopularTags)

	app.Listen(":" + strconv.Itoa(config.Port))
}
