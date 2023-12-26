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

type CommentService struct {
	DB         *sql.DB
	Config     *util.Config
	Validation *util.Validation
}

func (u *CommentService) getDB() *sql.DB {
	return u.DB
}

func (u *CommentService) getConfig() *util.Config {
	return u.Config
}

func (u *CommentService) getValidation() *util.Validation {
	return u.Validation
}

func (u *CommentService) CreateComment(c *fiber.Ctx) error {
	req := new(models.CommentCreation)
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
	queries := gen.New(u.getDB())
	if req.ContentType == "" {
		req.ContentType = models.CONTENT_TEXT
	}
	comment, err := queries.CreateComment(c.Context(), gen.CreateCommentParams{
		ParentID: sql.NullInt64{
			Int64: req.ParentId,
			Valid: true,
		},
		PosterID: req.PosterId,
		PostID:   req.PostId,
		Body:     req.Body,
		ContentType: gen.NullContentType{
			ContentType: gen.ContentType(req.ContentType),
			Valid:       req.ContentType != "",
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}

	tags, err := util.GetTags(c, req.Body, queries, nil, &comment)
	if req.ContentType == models.CONTENT_PICTURE && req.FileIds != nil {
		for _, id := range req.FileIds {
			_, err := queries.CreatePictureRelation(c.Context(), gen.CreatePictureRelationParams{
				PictureID: sql.NullInt64{
					Int64: id,
					Valid: true,
				},
				CommentID: sql.NullInt64{
					Int64: comment.ID,
					Valid: true,
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("Error while creating picture relations in db")
				return c.Status(fiber.StatusInternalServerError).SendStatus(500)
			}
		}
	}

	if req.ContentType == models.CONTENT_VIDEO && req.FileIds != nil {
		for _, id := range req.FileIds {
			_, err := queries.CreateVideoRelation(c.Context(), gen.CreateVideoRelationParams{
				VideoID: sql.NullInt64{
					Int64: id,
					Valid: true,
				},
				CommentID: sql.NullInt64{
					Int64: comment.ID,
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
		_, err = queries.CreateTagRelationForComment(c.Context(), gen.CreateTagRelationForCommentParams{
			TagID: sql.NullInt64{
				Int64: tag.ID,
				Valid: true,
			},
			CommentID: sql.NullInt64{
				Int64: comment.ID,
				Valid: true,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Error while creating using in db")
			return c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
	}

	tx.Commit()
	return c.JSON(models.Comment{
		Id:              comment.ID,
		PosterId:        comment.PosterID,
		ParentId:        comment.ParentID.Int64,
		PostId:          comment.PostID,
		Body:            comment.Body,
		CommentPictures: nil,
		CommentVideos:   nil,
		ContentType:     models.ContentType(comment.ContentType.ContentType),
		UpVotes:         0,
		DownVotes:       0,
		Created:         comment.Created.Time,
	})
}

func (u *CommentService) GetCommentById(c *fiber.Ctx) error {
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
	res, err := queries.GetComment(c.Context(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Send(nil)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}

	return c.JSON(getCommentNoVote(res.CommentsView))
}

func (u *CommentService) GetComments(c *fiber.Ctx) error {
	postId := c.QueryInt("postId", -1)
	commentId := c.QueryInt("commentId", -1)
	userId := c.QueryInt("userId", -1)
	sort := c.Query("sort", "latest")
	days := c.QueryInt("days", 7)
	if userId == -1 {
		return c.Status(fiber.StatusBadRequest).
			SendString("userId is required")
	}
	//TODO:get comment by popularity
	if postId == -1 && commentId == -1 {
		return c.Status(fiber.StatusBadRequest).
			SendString("Need to set either postId or commentId")
	}
	queries := gen.New(u.getDB())
	if postId != -1 {
		return getCommentsForPost(postId, sort, days, userId, queries, c)
	}
	res, err := queries.GetCommentsForComment(c.Context(), gen.GetCommentsForCommentParams{
		ID:     int64(commentId),
		UserID: int64(userId),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.SendStatus(fiber.StatusOK)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Comment
	for _, comment := range res {
		res, err := getComment(comment.CommentsView, comment.IsUpVote, comment.VoteType, comment.IsDeleted, queries, c)
		if err != nil {
			return err
		}
		list = append(list, *res)
	}
	return c.JSON(list)
}

func getCommentsForPost(postId int, sort string, days int, userId int, queries *gen.Queries, c *fiber.Ctx) error {
	if sort == "popular" {
		res, err := queries.GetCommentsForPostPopular(c.Context(), gen.GetCommentsForPostPopularParams{
			PostID: int64(postId),
			UserID: int64(userId),
			Days:   int32(days),
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.SendStatus(fiber.StatusOK)
			}
			log.Error().Err(err).Msg("Error while creating using in db")
			return c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		var list []models.Comment
		for _, comment := range res {
			res, err := getComment(comment.CommentsView, comment.IsUpVote, comment.VoteType, comment.IsDeleted, queries, c)
			if err != nil {
				return err
			}
			list = append(list, *res)
		}
		return c.JSON(list)
	}
	res, err := queries.GetCommentsForPostLatest(c.Context(), gen.GetCommentsForPostLatestParams{
		PostID: int64(postId),
		UserID: int64(userId),
		Days:   int32(days),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.SendStatus(fiber.StatusOK)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Comment
	for _, comment := range res {
		res, err := getComment(comment.CommentsView, comment.IsUpVote, comment.VoteType, comment.IsDeleted, queries, c)
		if err != nil {
			return err
		}
		list = append(list, *res)
	}
	return c.JSON(list)
}

func getComment(view gen.CommentsView, isUpVote sql.NullBool, voteType gen.NullVoteType, voteIsDeleted sql.NullBool, queries *gen.Queries, c *fiber.Ctx) (*models.Comment, error) {
	pictures, err := getPicturesForComment(view.ID, queries, c)
	videos, err := getVideosForComment(view.ID, queries, c)
	if err != nil {
		return nil, err
	}
	return &models.Comment{
		Id:              view.ID,
		PosterId:        view.PosterID,
		PostId:          view.PostID,
		ParentId:        view.ParentID.Int64,
		PosterName:      view.DisplayName,
		Body:            view.Body,
		ContentType:     models.ContentType(view.ContentType.ContentType),
		UpVotes:         view.UpVotes,
		DownVotes:       view.DownVotes,
		Created:         view.Created.Time,
		PosterPicture:   getPictureMeta(view.UserPicUrl, view.UserPicWidth, view.UserPicHeight, view.UserPicID.Int64),
		Vote:            getVote(isUpVote, voteType, voteIsDeleted),
		CommentPictures: pictures,
		CommentVideos:   videos,
	}, nil
}

func getCommentNoVote(view gen.CommentsView) models.Comment {
	return models.Comment{
		Id:            view.ID,
		PosterId:      view.PosterID,
		PostId:        view.PostID,
		ParentId:      view.ParentID.Int64,
		PosterName:    view.DisplayName,
		Body:          view.Body,
		ContentType:   models.ContentType(view.ContentType.ContentType),
		UpVotes:       view.UpVotes,
		DownVotes:     view.DownVotes,
		Created:       view.Created.Time,
		PosterPicture: getPictureMeta(view.UserPicUrl, view.UserPicWidth, view.UserPicHeight, view.UserPicID.Int64),
	}
}

func getPicturesForComment(commentId int64, queries *gen.Queries, c *fiber.Ctx) ([]*models.PictureMeta, error) {
	res, err := queries.GetPicturesForComment(c.Context(), sql.NullInt64{
		Int64: commentId,
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
	if len(slice) == 0 {
		return nil, nil
	}
	return slice, nil
}

func getVideosForComment(commentId int64, queries *gen.Queries, c *fiber.Ctx) ([]*models.VideoMeta, error) {
	res, err := queries.GetVideoForComment(c.Context(), sql.NullInt64{
		Int64: commentId,
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			var slice []*models.VideoMeta
			return slice, nil
		}
		log.Error().Err(err).Msg("Error getting pictures from db")
		return nil, c.SendStatus(fiber.StatusInternalServerError)
	}
	var slice []*models.VideoMeta
	for _, p := range res {
		slice = append(slice, getVideoMetaFromModel(p))
	}
	return slice, nil
}
