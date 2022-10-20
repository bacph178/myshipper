package dtos

import (
	"myshipper/models"
	"time"
)

func GetSummary(comment *models.Comment, includeUser, includeProduct bool) map[string]interface{} {
	result := map[string]interface{}{
		"id":         comment.ID,
		"content":    comment.Content,
		"created_at": comment.CreatedAt.UTC().Format(time.RFC1123),
		"updated_at": comment.UpdatedAt.UTC().Format(time.RFC1123),
	}
	if includeUser {
		result["user"] = map[string]interface{}{
			"id":       comment.User.ID,
			"username": comment.User.Username,
		}
	}
	if includeProduct {
		result["product"] = map[string]interface{}{
			"id":   comment.Product.ID,
			"name": comment.Product.Name,
			"slug": comment.Product.Slug,
		}
	}
	return result
}
