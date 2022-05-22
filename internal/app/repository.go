package app

import (
	"context"
	"hw2-tgbot/internal/models"
)

type repository interface {
	GetUsers(ctx context.Context) ([]*models.User, error)
	AddUser(ctx context.Context, user *models.User) error
	FindUser(ctx context.Context, tgID int64) (string, error)
}
