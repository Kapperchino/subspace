package server

import (
	"database/sql"
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
		Content: sql.NullString{
			String: req.Content,
			Valid:  true,
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
	return c.JSON(models.Comment{
		Id:          comment.ID,
		PosterId:    comment.PosterID,
		ParentId:    comment.ParentID.Int64,
		PostId:      comment.PostID,
		Body:        comment.Body,
		Content:     comment.Content.String,
		ContentType: models.ContentType(comment.ContentType.ContentType),
		UpVotes:     0,
		DownVotes:   0,
		Created:     comment.Created.Time,
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
		if err == sql.ErrNoRows {
			return c.Send(nil)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	return c.JSON(models.Comment{
		Id:          res.ID,
		PosterId:    res.PosterID,
		Body:        res.Body,
		ParentId:    res.ParentID.Int64,
		ContentType: models.ContentType(res.ContentType.ContentType),
		Content:     res.Content.String,
		UpVotes:     res.UpVotes,
		DownVotes:   res.DownVotes,
		Created:     res.Created.Time,
	})
}

func (u *CommentService) GetComments(c *fiber.Ctx) error {
	postId := c.QueryInt("postId", -1)
	commentId := c.QueryInt("commentId", -1)
	//TODO:get comment by popularity
	if postId == -1 && commentId == -1 {
		return c.Status(fiber.StatusBadRequest).
			SendString("Need to set either postId or commentId")
	}
	queries := gen.New(u.getDB())
	if postId != -1 {
		res, err := queries.GetCommentsForPost(c.Context(), int64(postId))
		if err != nil {
			if err == sql.ErrNoRows {
				return c.SendStatus(fiber.StatusOK)
			}
			log.Error().Err(err).Msg("Error while creating using in db")
			return c.Status(fiber.StatusInternalServerError).SendStatus(500)
		}
		var list []models.Comment
		for _, comment := range res {
			list = append(list, models.Comment{
				Id:          comment.ID,
				PosterId:    comment.PosterID,
				PostId:      comment.PostID,
				ParentId:    comment.ParentID.Int64,
				PosterName:  comment.DisplayName,
				Body:        comment.Body,
				Content:     comment.Content.String,
				ContentType: models.ContentType(comment.ContentType.ContentType),
				UpVotes:     comment.UpVotes,
				DownVotes:   comment.DownVotes,
				Created:     comment.Created.Time,
			})
		}
		return c.JSON(list)
	}
	res, err := queries.GetCommentsForComment(c.Context(), int64(commentId))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.SendStatus(fiber.StatusOK)
		}
		log.Error().Err(err).Msg("Error while creating using in db")
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	var list []models.Comment
	for _, comment := range res {
		list = append(list, models.Comment{
			Id:          comment.ID,
			PosterId:    comment.PosterID,
			PostId:      comment.PostID,
			ParentId:    comment.ParentID.Int64,
			Body:        comment.Body,
			Content:     comment.Content.String,
			ContentType: models.ContentType(comment.ContentType.ContentType),
			UpVotes:     comment.UpVotes,
			DownVotes:   comment.DownVotes,
			Created:     comment.Created.Time,
		})
	}
	return c.JSON(list)
}
