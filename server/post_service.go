package server

import (
	"database/sql"
	"github.com/Kapperchino/subspace/models"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/Kapperchino/subspace/util"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type PostService struct {
	DB           *sql.DB
	Config       *util.Config
	Validation   *util.Validation
	UploadClient *util.UploadClient
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
	queries := gen.New(u.getDB())
	if req.ContentType == "" {
		req.ContentType = models.CONTENT_TEXT
	}
	presigned, fileName, err := u.getPresigned(req.IsUpload, c)
	var content = ""
	if presigned != "" {
		content = "https://pub-cab547f3a0034c6083d1d10ab8298a3f.r2.dev/" + fileName
	} else {
		content = req.Content
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
			Valid:  req.Topic == "",
		},
		Body: sql.NullString{
			String: req.Body,
			Valid:  true,
		},
		Content: sql.NullString{
			String: content,
			Valid:  content != "",
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
	return c.JSON(models.Post{
		Id:          post.ID,
		SpaceId:     post.SpaceID.Int64,
		PosterId:    post.PosterID.Int64,
		Body:        post.Body.String,
		Topic:       post.Topic.String,
		Content:     post.Content.String,
		IsUpload:    req.IsUpload,
		Presigned:   presigned,
		ContentType: models.ContentType(post.ContentType.ContentType),
		UpVotes:     0,
		DownVotes:   0,
		Created:     post.Created.Time,
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
		if err == sql.ErrNoRows {
			return c.Send(nil)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var vote *models.Vote
	if res.ID_2.Valid {
		vote = &models.Vote{
			VoteId:          res.ID_2.Int64,
			UserId:          res.UserID.Int64,
			PostOrCommentId: res.PostOrCommentID.Int64,
			IsUpVote:        res.IsUpVote.Bool,
			VoteType:        models.VoteType(res.VoteType.VoteType),
			IsDeleted:       res.IsDeleted_2.Bool,
		}
	} else {
		vote = nil
	}
	return c.JSON(models.Post{
		Id:            res.ID,
		SpaceId:       res.SpaceID.Int64,
		PosterId:      res.PosterID.Int64,
		PosterName:    res.DisplayName,
		SpacePicture:  res.SpacePicture.String,
		Topic:         res.Topic.String,
		Body:          res.Body.String,
		ContentType:   models.ContentType(res.ContentType.ContentType),
		Content:       res.Content.String,
		UpVotes:       res.UpVotes,
		DownVotes:     res.DownVotes,
		Created:       res.Created.Time,
		Vote:          vote,
		SpaceParentId: res.ParentID,
		SpaceName:     res.SpaceName,
	})
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

func (u *PostService) getPresigned(isUpload bool, c *fiber.Ctx) (string, string, error) {
	if !isUpload {
		return "", "", nil
	}
	preSigned, key, err := u.UploadClient.Presign(c.Context())

	return preSigned.URL, key, err
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
			if err == sql.ErrNoRows {
				return nil, c.SendStatus(fiber.StatusOK)
			}
			log.Error().Err(err).Msg("Error while creating using in db")
			return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		var list []models.Post
		for _, post := range res {
			var vote *models.Vote
			if post.ID_2.Valid {
				vote = &models.Vote{
					VoteId:          post.ID_2.Int64,
					UserId:          post.UserID.Int64,
					PostOrCommentId: post.PostOrCommentID.Int64,
					IsUpVote:        post.IsUpVote.Bool,
					VoteType:        models.VoteType(post.VoteType.VoteType),
					IsDeleted:       post.IsDeleted_2.Bool,
				}
			} else {
				vote = nil
			}
			list = append(list, models.Post{
				Id:            post.ID,
				SpaceId:       post.SpaceID.Int64,
				PosterId:      post.PosterID.Int64,
				SpacePicture:  post.SpacePicture.String,
				Topic:         post.Topic.String,
				Content:       post.Content.String,
				PosterName:    post.DisplayName,
				ContentType:   models.ContentType(post.ContentType.ContentType),
				Body:          post.Body.String,
				UpVotes:       post.UpVotes,
				DownVotes:     post.DownVotes,
				Created:       post.Created.Time,
				Vote:          vote,
				SpaceParentId: post.ParentID,
				SpaceName:     post.SpaceName,
			})
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
		if err == sql.ErrNoRows {
			return nil, c.SendStatus(fiber.StatusOK)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Post
	for _, post := range res {
		var vote *models.Vote
		if post.ID_2.Valid {
			vote = &models.Vote{
				VoteId:          post.ID_2.Int64,
				UserId:          post.UserID.Int64,
				PostOrCommentId: post.PostOrCommentID.Int64,
				IsUpVote:        post.IsUpVote.Bool,
				VoteType:        models.VoteType(post.VoteType.VoteType),
				IsDeleted:       post.IsDeleted_2.Bool,
			}
		} else {
			vote = nil
		}

		list = append(list, models.Post{
			Id:            post.ID,
			SpaceId:       post.SpaceID.Int64,
			PosterId:      post.PosterID.Int64,
			SpacePicture:  post.SpacePicture.String,
			Topic:         post.Topic.String,
			Content:       post.Content.String,
			PosterName:    post.DisplayName,
			ContentType:   models.ContentType(post.ContentType.ContentType),
			Body:          post.Body.String,
			UpVotes:       post.UpVotes,
			DownVotes:     post.DownVotes,
			Created:       post.Created.Time,
			Vote:          vote,
			SpaceParentId: post.ParentID,
			SpaceName:     post.SpaceName,
		})
	}
	return list, nil
}

func (u *PostService) getPostsForSpaceByName(userId int64, parentId int64, spaceName string, isPopular bool, days int32, queries *gen.Queries, c *fiber.Ctx) ([]models.Post, error) {
	if !isPopular {
		res, err := queries.GetPostsForSpaceLatestByName(c.Context(), gen.GetPostsForSpaceLatestByNameParams{
			Name:     spaceName,
			ParentID: parentId,
			Days:     days,
			UserID:   userId,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, c.SendStatus(fiber.StatusOK)
			}
			log.Error().Err(err).Msg("Error while creating using in db")
			return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		var list []models.Post
		for _, post := range res {
			var vote *models.Vote
			if post.ID_2.Valid {
				vote = &models.Vote{
					VoteId:          post.ID_2.Int64,
					UserId:          post.UserID.Int64,
					PostOrCommentId: post.PostOrCommentID.Int64,
					IsUpVote:        post.IsUpVote.Bool,
					VoteType:        models.VoteType(post.VoteType.VoteType),
					IsDeleted:       post.IsDeleted_2.Bool,
				}
			} else {
				vote = nil
			}
			list = append(list, models.Post{
				Id:            post.ID,
				SpaceId:       post.SpaceID.Int64,
				PosterId:      post.PosterID.Int64,
				SpacePicture:  post.SpacePicture.String,
				Topic:         post.Topic.String,
				Content:       post.Content.String,
				PosterName:    post.DisplayName,
				ContentType:   models.ContentType(post.ContentType.ContentType),
				Body:          post.Body.String,
				UpVotes:       post.UpVotes,
				DownVotes:     post.DownVotes,
				Created:       post.Created.Time,
				Vote:          vote,
				SpaceParentId: post.ParentID,
				SpaceName:     post.SpaceName,
			})
		}
		return list, nil
	}
	res, err := queries.GetPostsForSpacePopularByName(c.Context(), gen.GetPostsForSpacePopularByNameParams{
		Name:     spaceName,
		ParentID: parentId,
		UserID:   userId,
		Days:     days,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, c.SendStatus(fiber.StatusOK)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Post
	for _, post := range res {
		var vote *models.Vote
		if post.ID_2.Valid {
			vote = &models.Vote{
				VoteId:          post.ID_2.Int64,
				UserId:          post.UserID.Int64,
				PostOrCommentId: post.PostOrCommentID.Int64,
				IsUpVote:        post.IsUpVote.Bool,
				VoteType:        models.VoteType(post.VoteType.VoteType),
				IsDeleted:       post.IsDeleted_2.Bool,
			}
		} else {
			vote = nil
		}

		list = append(list, models.Post{
			Id:            post.ID,
			SpaceId:       post.SpaceID.Int64,
			PosterId:      post.PosterID.Int64,
			SpacePicture:  post.SpacePicture.String,
			Topic:         post.Topic.String,
			Content:       post.Content.String,
			PosterName:    post.DisplayName,
			ContentType:   models.ContentType(post.ContentType.ContentType),
			Body:          post.Body.String,
			UpVotes:       post.UpVotes,
			DownVotes:     post.DownVotes,
			Created:       post.Created.Time,
			Vote:          vote,
			SpaceParentId: post.ParentID,
			SpaceName:     post.SpaceName,
		})
	}
	return list, nil
}

func (u *PostService) GetPostsForUser(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id", -1)
	if userId == -1 || err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	queries := gen.New(u.getDB())
	res, err := queries.GetPostsForUser(c.Context(), int64(userId))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.SendStatus(fiber.StatusOK)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Post
	for _, post := range res {
		var vote *models.Vote
		if post.ID_2.Valid {
			vote = &models.Vote{
				VoteId:          post.ID_2.Int64,
				UserId:          post.UserID.Int64,
				PostOrCommentId: post.PostOrCommentID.Int64,
				IsUpVote:        post.IsUpVote.Bool,
				VoteType:        models.VoteType(post.VoteType.VoteType),
				IsDeleted:       post.IsDeleted_2.Bool,
			}
		} else {
			vote = nil
		}
		list = append(list, models.Post{
			Id:            post.ID,
			SpaceId:       post.SpaceID.Int64,
			PosterId:      post.PosterID.Int64,
			SpacePicture:  post.SpacePicture.String,
			Topic:         post.Topic.String,
			Content:       post.Content.String,
			PosterName:    post.DisplayName,
			ContentType:   models.ContentType(post.ContentType.ContentType),
			Body:          post.Body.String,
			UpVotes:       post.UpVotes,
			DownVotes:     post.DownVotes,
			Created:       post.Created.Time,
			Vote:          vote,
			SpaceParentId: post.ParentID,
			SpaceName:     post.SpaceName,
		})
	}
	return c.JSON(list)
}

func (u *PostService) getPostsForHome(userId int64, isPopular bool, days int32, queries *gen.Queries, c *fiber.Ctx) ([]models.Post, error) {
	if !isPopular {
		res, err := queries.GetPostsForHomeLatest(c.Context(), gen.GetPostsForHomeLatestParams{
			Days:   days,
			UserID: userId,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, c.SendStatus(fiber.StatusOK)
			}
			log.Error().Err(err).Msg("Error while creating using in db")
			return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		var list []models.Post
		for _, post := range res {
			var vote *models.Vote
			if post.ID_2.Valid {
				vote = &models.Vote{
					VoteId:          post.ID_2.Int64,
					UserId:          post.UserID.Int64,
					PostOrCommentId: post.PostOrCommentID.Int64,
					IsUpVote:        post.IsUpVote.Bool,
					VoteType:        models.VoteType(post.VoteType.VoteType),
					IsDeleted:       post.IsDeleted_2.Bool,
				}
			} else {
				vote = nil
			}
			list = append(list, models.Post{
				Id:            post.ID,
				SpaceId:       post.SpaceID.Int64,
				PosterId:      post.PosterID.Int64,
				SpacePicture:  post.SpacePicture.String,
				Topic:         post.Topic.String,
				Content:       post.Content.String,
				PosterName:    post.DisplayName,
				ContentType:   models.ContentType(post.ContentType.ContentType),
				Body:          post.Body.String,
				UpVotes:       post.UpVotes,
				DownVotes:     post.DownVotes,
				Created:       post.Created.Time,
				Vote:          vote,
				SpaceParentId: post.ParentID,
				SpaceName:     post.SpaceName,
			})
		}
		return list, nil
	}
	res, err := queries.GetPostsForHomePopular(c.Context(), gen.GetPostsForHomePopularParams{
		UserID: userId,
		Days:   days,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, c.SendStatus(fiber.StatusOK)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Post
	for _, post := range res {
		var vote *models.Vote
		if post.ID_2.Valid {
			vote = &models.Vote{
				VoteId:          post.ID_2.Int64,
				UserId:          post.UserID.Int64,
				PostOrCommentId: post.PostOrCommentID.Int64,
				IsUpVote:        post.IsUpVote.Bool,
				VoteType:        models.VoteType(post.VoteType.VoteType),
				IsDeleted:       post.IsDeleted_2.Bool,
			}
		} else {
			vote = nil
		}

		list = append(list, models.Post{
			Id:            post.ID,
			SpaceId:       post.SpaceID.Int64,
			PosterId:      post.PosterID.Int64,
			SpacePicture:  post.SpacePicture.String,
			Topic:         post.Topic.String,
			Content:       post.Content.String,
			PosterName:    post.DisplayName,
			ContentType:   models.ContentType(post.ContentType.ContentType),
			Body:          post.Body.String,
			UpVotes:       post.UpVotes,
			DownVotes:     post.DownVotes,
			Created:       post.Created.Time,
			Vote:          vote,
			SpaceParentId: post.ParentID,
			SpaceName:     post.SpaceName,
		})
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
			if err == sql.ErrNoRows {
				return nil, c.SendStatus(fiber.StatusOK)
			}
			log.Error().Err(err).Msg("Error while creating using in db")
			return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		var list []models.Post
		for _, post := range res {
			var vote *models.Vote
			if post.ID_2.Valid {
				vote = &models.Vote{
					VoteId:          post.ID_2.Int64,
					UserId:          post.UserID.Int64,
					PostOrCommentId: post.PostOrCommentID.Int64,
					IsUpVote:        post.IsUpVote.Bool,
					VoteType:        models.VoteType(post.VoteType.VoteType),
					IsDeleted:       post.IsDeleted_2.Bool,
				}
			} else {
				vote = nil
			}
			list = append(list, models.Post{
				Id:            post.ID,
				SpaceId:       post.SpaceID.Int64,
				PosterId:      post.PosterID.Int64,
				SpacePicture:  post.SpacePicture.String,
				Topic:         post.Topic.String,
				Content:       post.Content.String,
				PosterName:    post.DisplayName,
				ContentType:   models.ContentType(post.ContentType.ContentType),
				Body:          post.Body.String,
				UpVotes:       post.UpVotes,
				DownVotes:     post.DownVotes,
				Created:       post.Created.Time,
				Vote:          vote,
				SpaceParentId: post.ParentID,
				SpaceName:     post.SpaceName,
			})
		}
		return list, nil
	}
	res, err := queries.GetPostsForUserSubscriptionPopular(c.Context(), gen.GetPostsForUserSubscriptionPopularParams{
		UserID: userId,
		Days:   days,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, c.SendStatus(fiber.StatusOK)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Post
	for _, post := range res {
		var vote *models.Vote
		if post.ID_2.Valid {
			vote = &models.Vote{
				VoteId:          post.ID_2.Int64,
				UserId:          post.UserID.Int64,
				PostOrCommentId: post.PostOrCommentID.Int64,
				IsUpVote:        post.IsUpVote.Bool,
				VoteType:        models.VoteType(post.VoteType.VoteType),
				IsDeleted:       post.IsDeleted_2.Bool,
			}
		} else {
			vote = nil
		}

		list = append(list, models.Post{
			Id:            post.ID,
			SpaceId:       post.SpaceID.Int64,
			PosterId:      post.PosterID.Int64,
			SpacePicture:  post.SpacePicture.String,
			Topic:         post.Topic.String,
			Content:       post.Content.String,
			PosterName:    post.DisplayName,
			ContentType:   models.ContentType(post.ContentType.ContentType),
			Body:          post.Body.String,
			UpVotes:       post.UpVotes,
			DownVotes:     post.DownVotes,
			Created:       post.Created.Time,
			Vote:          vote,
			SpaceParentId: post.ParentID,
			SpaceName:     post.SpaceName,
		})
	}
	return list, nil
}
