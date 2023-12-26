package util

import (
	"database/sql"
	"errors"
	gen "github.com/Kapperchino/subspace/sql/generated"
	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"strings"
)

func GetTags(c *fiber.Ctx, body string, queries *gen.Queries, post *gen.Post, comment *gen.Comment) ([]gen.Tag, error) {
	var tags []gen.Tag
	if body != "" {
		p := strings.NewReader(body)
		doc, _ := goquery.NewDocumentFromReader(p)
		var mentions []string
		var hashTags []string

		doc.Find("span.mentions").Each(func(i int, s *goquery.Selection) {
			mention, valid := s.Attr("data-id")
			if valid {
				mentions = append(mentions, mention)
			}
		})

		doc.Find("span.hashtags").Each(func(i int, s *goquery.Selection) {
			hashTag, valid := s.Attr("data-id")
			if valid {
				hashTags = append(hashTags, hashTag)
			}
		})

		for _, s := range mentions {
			user, err := queries.GetUserFromAddress(c.Context(), sql.NullString{
				String: s,
				Valid:  true,
			})
			if err != nil && errors.Is(err, sql.ErrNoRows) {
				return nil, c.SendStatus(fiber.StatusNotFound)
			}
			if err != nil {
				log.Error().Err(err).Msg("Error while creating using in db")
				return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
			}

			if post != nil {
				_, err = queries.CreateUserAddress(c.Context(), gen.CreateUserAddressParams{
					FromUserID: sql.NullInt64{Int64: post.PosterID.Int64, Valid: true},
					ToUserID:   sql.NullInt64{Int64: user.User.ID, Valid: true},
					PostID:     sql.NullInt64{Int64: post.ID, Valid: true},
				})
			}

			if comment != nil {
				_, err = queries.CreateUserAddress(c.Context(), gen.CreateUserAddressParams{
					FromUserID: sql.NullInt64{Int64: comment.PosterID, Valid: true},
					ToUserID:   sql.NullInt64{Int64: user.User.ID, Valid: true},
					PostID:     sql.NullInt64{Int64: post.ID, Valid: true},
				})
			}

			if err != nil {
				log.Error().Err(err).Msg("Error while creating using in db")
				return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
			}
		}

		for _, s := range hashTags {
			tag, err := queries.GetTag(c.Context(), s)
			if err != nil && errors.Is(err, sql.ErrNoRows) || tag.ID == 0 {
				tag, err = queries.CreateTag(c.Context(), s)
				if err != nil {
					log.Error().Err(err).Msg("Error while creating using in db")
					return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
				}
			}
			if err != nil {
				log.Error().Err(err).Msg("Error while creating using in db")
				return nil, c.Status(fiber.StatusInternalServerError).SendStatus(500)
			}
			tags = append(tags, tag)
		}
	}
	return tags, nil
}
