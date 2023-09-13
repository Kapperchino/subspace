package server

import (
	"database/sql"
	"errors"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/Kapperchino/subspace/util"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"strings"
)

type PostService struct {
	DB         *sql.DB
	Config     *util.Config
	Validation *util.Validation
}

func (u *PostService) getDB() *sql.DB {
	return u.DB
}

func (u *PostService) getConfig() *util.Config {
	return u.Config
}

func (u *PostService) getValidation() *util.Validation {
	return u.Validation
}

func (u *PostService) CreatePost(c *fiber.Ctx) error {
	req := new(models.PostCreation)
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
	if req.ContentType == "" {
		req.ContentType = models.CONTENT_TEXT
	}

	post, err := queries.CreatePost(c.Context(), gen.CreatePostParams{
		SpaceID: sql.NullInt64{
			Int64: req.SpaceId,
			Valid: true,
		},
		PosterID: sql.NullInt64{
			Int64: req.PosterId,
			Valid: true,
		},
		Topic: sql.NullString{
			String: req.Topic,
			Valid:  req.Topic != "",
		},
		Body: sql.NullString{
			String: req.Body,
			Valid:  req.Body != "",
		},
		Link: sql.NullString{
			String: req.Link,
			Valid:  req.Link != "",
		},
		ContentType: gen.NullContentType{
			ContentType: gen.ContentType(req.ContentType),
			Valid:       req.ContentType != "",
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}

	var tags []gen.Tag
	if req.Body != "" {
		r := util.TagRegex
		res := r.FindAllString(req.Body, -1)
		if res != nil {
			for _, s := range res {
				trimmed := strings.TrimSpace(s)
				switch trimmed[0] {
				case '@':
					address := trimmed[1:]
					user, err := queries.GetUserFromAddress(c.Context(), sql.NullString{
						String: address,
						Valid:  true,
					})
					if err != nil && errors.Is(err, sql.ErrNoRows) {
						return c.SendStatus(fiber.StatusNotFound)
					}
					if err != nil {
						log.Error().Err(err).Msg("Error while creating using in db")
						return c.Status(fiber.StatusInternalServerError).SendStatus(500)
					}
					_, err = queries.CreateUserAddress(c.Context(), gen.CreateUserAddressParams{
						FromUserID: sql.NullInt64{Int64: req.PosterId, Valid: true},
						ToUserID:   sql.NullInt64{Int64: user.User.ID, Valid: true},
						PostID:     sql.NullInt64{Int64: post.ID, Valid: true},
					})
					if err != nil {
						log.Error().Err(err).Msg("Error while creating using in db")
						return c.Status(fiber.StatusInternalServerError).SendStatus(500)
					}
					break
				case '#':
					tag, err := queries.GetTag(c.Context(), trimmed[1:])
					if err != nil && errors.Is(err, sql.ErrNoRows) || tag.ID == 0 {
						tag, err = queries.CreateTag(c.Context(), trimmed[1:])
						if err != nil {
							log.Error().Err(err).Msg("Error while creating using in db")
							return c.Status(fiber.StatusInternalServerError).SendStatus(500)
						}
					}
					if err != nil {
						log.Error().Err(err).Msg("Error while creating using in db")
						return c.Status(fiber.StatusInternalServerError).SendStatus(500)
					}
					tags = append(tags, tag)
					break
				}
			}
		}
	}

	if req.ContentType == models.CONTENT_PICTURE && req.FileIds != nil {
		for _, id := range req.FileIds {
			_, err := queries.CreatePictureRelation(c.Context(), gen.CreatePictureRelationParams{
				PictureID: sql.NullInt64{
					Int64: id,
					Valid: true,
				},
				PostID: sql.NullInt64{
					Int64: post.ID,
					Valid: true,
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("Error while creating picture relations in db")
				return c.Status(fiber.StatusInternalServerError).SendStatus(500)
			}
		}
	}

	for _, tag := range tags {
		_, err = queries.CreateTagRelationForPost(c.Context(), gen.CreateTagRelationForPostParams{
			TagID: sql.NullInt64{
				Int64: tag.ID,
				Valid: true,
			},
			PostID: sql.NullInt64{
				Int64: post.ID,
				Valid: true,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Error while creating using in db")
			return c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.JSON(models.Post{
		Id:           post.ID,
		SpaceId:      post.SpaceID.Int64,
		PosterId:     post.PosterID.Int64,
		Body:         post.Body.String,
		Topic:        post.Topic.String,
		Link:         post.Link.String,
		PostPictures: nil,
		ContentType:  models.ContentType(post.ContentType.ContentType),
		UpVotes:      0,
		DownVotes:    0,
		Created:      post.Created.Time,
	})
}

func (u *PostService) GetPostById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id", -1)
	userId := c.QueryInt("userId", -1)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if id == -1 || userId == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(u.getDB())
	res, err := queries.GetPost(c.Context(), gen.GetPostParams{
		ID:     int64(id),
		UserID: int64(userId),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Send(nil)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	postModel, err := getPost(res.PostsView, res.IsUpVote, res.VoteType, queries, c)
	if err != nil {
		return err
	}
	return c.JSON(postModel)
}

func (u *PostService) GetPostsForTag(c *fiber.Ctx) error {
	tag := c.Params("name")
	userId := c.QueryInt("userId", -1)
	sort := c.Query("sort", "latest")
	days := c.QueryInt("days", 7)
	if tag == "" || userId == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if sort != "latest" && sort != "popular" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if days > 365 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(u.getDB())

	list, err := getPostsForTag(tag, int64(userId), sort == "popular", int32(days), queries, c)
	if err != nil {
		return err
	}
	return c.JSON(list)
}

func (u *PostService) GetPostsByName(c *fiber.Ctx) error {
	spaceName := c.Query("space", "")
	parentId := c.QueryInt("parentId", -1)
	userId := c.QueryInt("userId", -1)
	sort := c.Query("sort", "latest")
	days := c.QueryInt("days", 7)
	if spaceName == "" || userId == -1 || parentId == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if sort != "latest" && sort != "popular" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if days > 365 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(u.getDB())

	list, err := u.getPostsForSpaceByName(int64(userId), int64(parentId), spaceName, sort == "popular", int32(days), queries, c)
	if err != nil {
		return err
	}
	return c.JSON(list)
}

func (u *PostService) GetPosts(c *fiber.Ctx) error {
	spaceId, _ := c.ParamsInt("id", -1)
	userId := c.QueryInt("userId", -1)
	sort := c.Query("sort", "latest")
	days := c.QueryInt("days", 7)
	if spaceId == -1 || userId == -1 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if sort != "latest" && sort != "popular" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if days > 365 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if spaceId == -1 {
		spaceId = 1
	}
	queries := gen.New(u.getDB())

	list, err := u.getPosts(int64(userId), int64(spaceId), sort == "popular", int32(days), queries, c)
	if err != nil {
		return err
	}
	return c.JSON(list)
}

func (u *PostService) GetPostsForSubscription(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id", -1)
	sort := c.Query("sort", "latest")
	days := c.QueryInt("days", 7)
	if userId == -1 || err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if sort != "latest" && sort != "popular" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if days > 365 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(u.getDB())
	list, err := u.getPostsForUserSubscription(int64(userId), sort == "popular", int32(days), queries, c)
	if err != nil {
		return err
	}
	return c.JSON(list)
}

func (u *PostService) getPosts(userId int64, spaceId int64, isPopular bool, days int32, queries *gen.Queries, c *fiber.Ctx) ([]models.Post, error) {
	if spaceId == 1 {
		return u.getPostsForHome(userId, isPopular, days, queries, c)
	}
	if !isPopular {
		res, err := queries.GetPostsForSpaceLatest(c.Context(), gen.GetPostsForSpaceLatestParams{
			SpaceID: sql.NullInt64{
				Int64: spaceId,
				Valid: true,
			},
			Days:   days,
			UserID: userId,
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
			postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, queries, c)
			if err != nil {
				return nil, err
			}
			list = append(list, *postModel)
		}
		return list, nil
	}
	res, err := queries.GetPostsForSpacePopular(c.Context(), gen.GetPostsForSpacePopularParams{
		SpaceID: sql.NullInt64{
			Int64: spaceId,
			Valid: true,
		},
		UserID: userId,
		Days:   days,
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
		postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, queries, c)
		if err != nil {
			return nil, err
		}
		list = append(list, *postModel)
	}
	return list, nil
}

func (u *PostService) getPostsForSpaceByName(userId int64, parentId int64, spaceName string, isPopular bool, days int32, queries *gen.Queries, c *fiber.Ctx) ([]models.Post, error) {
	if !isPopular {
		res, err := queries.GetPostsForSpaceLatestByName(c.Context(), gen.GetPostsForSpaceLatestByNameParams{
			SpaceName: spaceName,
			ParentID:  parentId,
			Days:      days,
			UserID:    userId,
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
			postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, queries, c)
			if err != nil {
				return nil, err
			}
			list = append(list, *postModel)
		}
		return list, nil
	}
	res, err := queries.GetPostsForSpacePopularByName(c.Context(), gen.GetPostsForSpacePopularByNameParams{
		SpaceName: spaceName,
		ParentID:  parentId,
		UserID:    userId,
		Days:      days,
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
		postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, queries, c)
		if err != nil {
			return nil, err
		}
		list = append(list, *postModel)
	}
	return list, nil
}

func (u *PostService) GetPostsForUser(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id", -1)
	sort := c.Query("sort", "latest")
	days := c.QueryInt("days", 7)
	if userId == -1 || err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if sort != "latest" && sort != "popular" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if days > 365 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(u.getDB())
	list, err := u.getPostsForUser(int64(userId), sort == "popular", int32(days), queries, c)
	if err != nil {
		return err
	}
	return c.JSON(list)
}

func (u *PostService) getPostsForUser(userId int64, isPopular bool, days int32, queries *gen.Queries, c *fiber.Ctx) ([]models.Post, error) {
	if !isPopular {
		res, err := queries.GetPostsForUserLatest(c.Context(), gen.GetPostsForUserLatestParams{
			Days:   days,
			UserID: userId,
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
			postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, queries, c)
			if err != nil {
				return nil, err
			}
			list = append(list, *postModel)
		}
		return list, nil
	}
	res, err := queries.GetPostsForUserPopular(c.Context(), gen.GetPostsForUserPopularParams{
		UserID: userId,
		Days:   days,
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
		postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, queries, c)
		if err != nil {
			return nil, err
		}
		list = append(list, *postModel)
	}
	return list, nil
}

func (u *PostService) getPostsForHome(userId int64, isPopular bool, days int32, queries *gen.Queries, c *fiber.Ctx) ([]models.Post, error) {
	if !isPopular {
		res, err := queries.GetPostsForHomeLatest(c.Context(), gen.GetPostsForHomeLatestParams{
			Days:   days,
			UserID: userId,
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
			postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, queries, c)
			if err != nil {
				return nil, err
			}
			list = append(list, *postModel)
		}
		return list, nil
	}
	res, err := queries.GetPostsForHomePopular(c.Context(), gen.GetPostsForHomePopularParams{
		UserID: userId,
		Days:   days,
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
		postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, queries, c)
		if err != nil {
			return nil, err
		}
		list = append(list, *postModel)
	}
	return list, nil
}

func (u *PostService) getPostsForUserSubscription(userId int64, isPopular bool, days int32, queries *gen.Queries, c *fiber.Ctx) ([]models.Post, error) {
	if !isPopular {
		res, err := queries.GetPostsForUserSubscriptionLatest(c.Context(), gen.GetPostsForUserSubscriptionLatestParams{
			Days:   days,
			UserID: userId,
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
			postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, queries, c)
			if err != nil {
				return nil, err
			}
			list = append(list, *postModel)
		}
		return list, nil
	}
	res, err := queries.GetPostsForUserSubscriptionPopular(c.Context(), gen.GetPostsForUserSubscriptionPopularParams{
		UserID: userId,
		Days:   days,
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
		postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, queries, c)
		if err != nil {
			return nil, err
		}
		list = append(list, *postModel)
	}
	return list, nil
}

func getPostsForTag(tag string, userId int64, isPopular bool, days int32, queries *gen.Queries, c *fiber.Ctx) ([]models.Post, error) {
	if !isPopular {
		res, err := queries.GetPostsWithTagsLatest(c.Context(), gen.GetPostsWithTagsLatestParams{
			Days:   days,
			UserID: userId,
			Name:   tag,
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
			postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, queries, c)
			if err != nil {
				return nil, err
			}
			list = append(list, *postModel)
		}
		return list, nil
	}
	res, err := queries.GetPostsWithTagsPopular(c.Context(), gen.GetPostsWithTagsPopularParams{
		Days:   days,
		UserID: userId,
		Name:   tag,
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
		postModel, err := getPost(post.PostsView, post.IsUpVote, post.VoteType, queries, c)
		if err != nil {
			return nil, err
		}
		list = append(list, *postModel)
	}
	return list, nil
}

func getPictureMeta(url sql.NullString, width sql.NullInt64, height sql.NullInt64, id int64) *models.PictureMeta {
	if !url.Valid {
		return nil
	}
	return &models.PictureMeta{
		Url:    url.String,
		Width:  width.Int64,
		Height: height.Int64,
		Id:     id,
	}
}

func getVote(isUpVote sql.NullBool, voteType gen.NullVoteType) *models.Vote {
	if !isUpVote.Valid {
		return nil
	}
	return &models.Vote{
		IsUpVote: isUpVote.Bool,
		VoteType: models.VoteType(voteType.VoteType),
	}
}

func getPictureMetaFromModel(picture gen.Picture) *models.PictureMeta {
	if picture.ID == 0 {
		return nil
	}
	return &models.PictureMeta{
		Url:    picture.Url,
		Width:  picture.Width,
		Height: picture.Height,
		Id:     picture.ID,
	}
}

func getPicturesForPost(postId int64, queries *gen.Queries, c *fiber.Ctx) ([]*models.PictureMeta, error) {
	res, err := queries.GetPicturesForPost(c.Context(), sql.NullInt64{
		Int64: postId,
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			var slice []*models.PictureMeta
			return slice, nil
		}
		log.Error().Err(err).Msg("Error getting pictures from db")
		return nil, c.SendStatus(fiber.StatusInternalServerError)
	}
	var slice []*models.PictureMeta
	for _, p := range res {
		slice = append(slice, getPictureMetaFromModel(p))
	}
	return slice, nil
}

func getPost(postView gen.PostsView, isUpvote sql.NullBool, voteType gen.NullVoteType, queries *gen.Queries, c *fiber.Ctx) (*models.Post, error) {
	pictures, err := getPicturesForPost(postView.ID, queries, c)
	if err != nil {
		return nil, err
	}
	return &models.Post{
		Id:            postView.ID,
		SpaceId:       postView.SpaceID.Int64,
		PosterId:      postView.PosterID.Int64,
		Topic:         postView.Topic.String,
		PosterName:    postView.DisplayName,
		ContentType:   models.ContentType(postView.ContentType.ContentType),
		Body:          postView.Body.String,
		Link:          postView.Link.String,
		UpVotes:       postView.UpVotes,
		DownVotes:     postView.DownVotes,
		PosterPicture: getPictureMeta(postView.UserPicUrl, postView.UserPicWidth, postView.UserPicHeight, postView.UserPicID.Int64),
		SpacePicture:  getPictureMeta(postView.SpaceSmallPicUrl, postView.SpaceSmallPicWidth, postView.SpaceSmallPicHeight, postView.SpaceSmallPicID.Int64),
		PostPictures:  pictures,
		Created:       postView.Created.Time,
		Vote:          getVote(isUpvote, voteType),
		SpaceParentId: postView.ParentID,
		SpaceName:     postView.SpaceName,
		CommentsCount: postView.CommentCount,
	}, nil
}
