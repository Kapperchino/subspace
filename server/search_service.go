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

type SearchService struct {
	DB         *sql.DB
	Config     *util.Config
	Validation *util.Validation
}

func (u *SearchService) getDB() *sql.DB {
	return u.DB
}

func (u *SearchService) getConfig() *util.Config {
	return u.Config
}

func (u *SearchService) getValidation() *util.Validation {
	return u.Validation
}

func (u *SearchService) SearchSpace(c *fiber.Ctx) error {
	search := c.Query("term")
	prefix := c.Query("prefix")
	// get all spaces
	if search == "" && prefix == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(u.getDB())
	if prefix != "" {
		res, err := queries.SpacePrefixSearch(c.Context(), prefix)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.Send(nil)
			}
			log.Error().Err(err).Msg("Error while searching db")
			return c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		var list []models.SpacePrefixRes
		for _, spaceView := range res {
			space := spaceView.SpacesView
			list = append(list, models.SpacePrefixRes{
				ID:           space.ID,
				ParentID:     space.ParentID,
				Name:         space.Name,
				SmallPicture: getPictureMeta(space.SpaceSmallPicUrl, space.SpaceSmallPicWidth, space.SpaceSmallPicHeight, space.SpaceSmallPictureID.Int64),
				SubCount:     space.SubCount,
			})
		}
		return c.JSON(list)
	}
	res, err := queries.SearchSpace(c.Context(), search)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Send(nil)
		}
		log.Error().Err(err).Msg("Error while searching db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Space
	for _, space := range res {
		list = append(list, models.Space{
			ID:                space.ID,
			ParentID:          space.ParentID,
			Name:              space.Name,
			Description:       space.Description.String,
			SmallPicture:      getPictureMeta(space.SpaceSmallPicUrl, space.SpaceSmallPicWidth, space.SpaceSmallPicHeight, space.SpaceSmallPictureID.Int64),
			BackgroundPicture: getPictureMeta(space.SpaceBackgroundPictureUrl, space.SpaceBackgroundPictureWidth, space.SpaceBackgroundPictureHeight, space.SpaceSmallPictureID.Int64),
			SubCount:          space.SubCount,
		})
	}
	return c.JSON(list)
}

func (u *SearchService) SearchPosts(c *fiber.Ctx) error {
	search := c.Query("term")
	userId := c.QueryInt("userId", -1)
	days := c.QueryInt("days", 7)
	isTag := c.QueryBool("isTag", false)
	sort := c.Query("sort", "latest")
	start := c.QueryInt("start", 0)
	// get all spaces
	if search == "" || userId == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(u.getDB())
	if isTag {
		list, err := getPostsForTag(search, int64(userId), sort == "popular", int32(days), queries, c, int32(start))
		if err != nil {
			log.Error().Err(err).Msg("Error getting posts for tag")
			return err
		}
		return c.JSON(list)
	}
	list, err := searchPosts(search, int64(userId), sort == "popular", int32(days), queries, c)
	if err != nil {
		log.Error().Err(err).Msgf("Error search term %s", search)
		return err
	}
	return c.JSON(list)
}

func searchPosts(term string, userId int64, isPopular bool, days int32, queries *gen.Queries, c *fiber.Ctx) ([]models.Post, error) {
	if !isPopular {
		res, err := queries.SearchPostLatest(c.Context(), gen.SearchPostLatestParams{
			PlaintoTsquery: term,
			Days:           days,
			UserID:         userId,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, c.SendStatus(fiber.StatusOK)
			}
			log.Error().Err(err).Msg("Error while creating using in db")
			return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		var list []models.Post
		for _, post := range res {
			postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, post.IsDeleted, queries, c)
			if err != nil {
				return nil, err
			}
			list = append(list, *postModel)
		}
		return list, nil
	}
	res, err := queries.SearchPostPopular(c.Context(), gen.SearchPostPopularParams{
		PlaintoTsquery: term,
		Days:           days,
		UserID:         userId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, c.SendStatus(fiber.StatusOK)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Post
	for _, post := range res {
		postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, post.IsDeleted, queries, c)
		if err != nil {
			return nil, err
		}
		list = append(list, *postModel)
	}
	return list, nil
}
